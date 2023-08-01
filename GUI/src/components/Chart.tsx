import Chart, { ChartOptions } from "chart.js/auto";
import "chartjs-adapter-moment";
import React, { useEffect, useRef, useState } from "react";

interface IDataPoint {
  Time: number;
  Value: number;
}

interface IData {
  Sensor: string;
  SetPoint: number;
  SetPointTs: number;
  Value: number;
  Ts: number;
  Error: null;
  TsSent: number;
}

interface ChartJsCanvasProps {
  name: string;
  datapoint: IDataPoint | undefined;
}

const ChartJsCanvas: React.FC<ChartJsCanvasProps> = ({ name, datapoint }) => {
  const chartRef = useRef<HTMLCanvasElement>(null);
  const chartInstanceRef = useRef<Chart>();

  useEffect(() => {
    if (chartRef.current) {
      const chartData = {
        labels: [],
        datasets: [
          {
            label: name,
            data: [],
            borderColor: "rgba(255, 99, 132, 1)",
            borderWidth: 1,
            fill: false,
          },
        ],
      };
      const chartOptions: ChartOptions = {
        responsive: true,
        maintainAspectRatio: false,
        animation: false,
        plugins: {
          tooltip: {
            enabled: false,
          },
        },
        scales: {
          x: {
            type: "time",
            time: {
              round: "millisecond",
              displayFormats: { seconds: "HH:mm:SS" },
            },
          },
          y: {},
        },
      };
      chartInstanceRef.current = new Chart(chartRef.current, {
        type: "line",
        data: chartData,
        options: chartOptions,
      });
    }
    return () => {
      chartInstanceRef.current?.destroy();
    };
  }, [name]);

  useEffect(() => {
    if (chartInstanceRef.current) {
      if (datapoint) {
        // add label X a
        chartInstanceRef.current.data?.labels?.push(datapoint.Time);
        // add data y
        chartInstanceRef.current.data.datasets[0].data.push(datapoint.Value);
      }
    }
  }, [datapoint]);

  useEffect(() => {
    const intervalId = setInterval(() => {
      chartInstanceRef?.current?.update();
    }, 200);

    return () => clearInterval(intervalId);
  }, []);

  return <canvas ref={chartRef}></canvas>;
};

interface JsChartProps {
  URL: string;
}

const ChartContainer: React.FC<JsChartProps> = ({ URL }) => {
  const [pdata, setpData] = useState<IDataPoint>();
  const [tdata, settData] = useState<IDataPoint>();

  useEffect(() => {
    const socket = new WebSocket(URL);

    socket.addEventListener("message", (event) => {
      const newDataPoint: IData = JSON.parse(event.data);
      if (newDataPoint.Sensor === "Pressure1") {
        setpData({
          Time: newDataPoint.Ts / 1000000,
          Value: newDataPoint.Value,
        });
      } else if (newDataPoint.Sensor === "Temperature1") {
        settData({
          Time: newDataPoint.Ts / 1000000,
          Value: newDataPoint.Value,
        });
      } else {
        console.log("Unknown sensor");
      }
    });

    return () => {
      // Close the WebSocket connection when the component unmounts
      // Make sure the connection is only closed if it's ready
      if (socket.readyState === 1) {
        console.log("Closing socket");
        socket.close();
      }
    };
  }, [URL]);

  return (
    <div className="ChartsContainer">
      <div className="chart">
        <ChartJsCanvas name={"Temperature"} datapoint={tdata} />
      </div>
      <div className="chart">
        <ChartJsCanvas name={"Pressure"} datapoint={pdata} />
      </div>
    </div>
  );
};

export default ChartContainer;
