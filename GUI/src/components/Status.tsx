import React, { useEffect, useState } from "react";
import SensorStatus, { SensorData } from "./SensorStatus";
import Led from "./StatusLigth";
import { IData } from "./types";

interface StatusProps {
  URL: string;
  started: boolean;
}

const StatusBar: React.FC<StatusProps> = ({ URL, started }) => {
  const [pdata, setpData] = useState<SensorData>();
  const [tdata, settData] = useState<SensorData>();

  useEffect(() => {
    const socket = new WebSocket(URL);

    socket.addEventListener("message", (event) => {
      const newDataPoint: IData = JSON.parse(event.data);
      if (newDataPoint.Sensor === "Pressure1") {
        setpData({
          Sensor: newDataPoint.Sensor,
          SetPoint: newDataPoint.SetPoint,
          Value: newDataPoint.Value,
          Error: newDataPoint.Error,
        });
      } else if (newDataPoint.Sensor === "Temperature1") {
        settData({
          Sensor: newDataPoint.Sensor,
          SetPoint: newDataPoint.SetPoint,
          Value: newDataPoint.Value,
          Error: newDataPoint.Error,
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
    <div className="status-bar">
      {started ? (
        <>
          <SensorStatus sensor={pdata} />
          <SensorStatus sensor={tdata} />
        </>
      ) : (
        <div> Start the oven first </div>
      )}
      <Led color="green" size={5} on={started} />
    </div>
  );
};

export default StatusBar;
