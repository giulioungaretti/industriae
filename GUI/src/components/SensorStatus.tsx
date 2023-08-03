export interface SensorData {
  Sensor: string;
  SetPoint: number;
  Value: Number;
  Error: null;
}

interface SensorStatusProps {
  sensor: SensorData | undefined;
}

const SensorStatus: React.FC<SensorStatusProps> = ({ sensor }) => {
  let statusClass = "";
  let indicatorClass = "";

  if (sensor) {
    if (sensor.Error) {
      statusClass = "sensor-status__status--error";
      indicatorClass = "sensor-status__indicator--error";
    } else if (sensor.SetPoint !== sensor.Value) {
      statusClass = "sensor-status__status--warning";
      indicatorClass = "sensor-status__indicator--warning";
    } else {
      statusClass = "sensor-status__status--ok";
      indicatorClass = "sensor-status__indicator--ok";
    }
  }
  return (
    <div className="sensor-status">
      {sensor && (
        <>
          <div className={`sensor-status__indicator ${indicatorClass}`}></div>
          <div className={`sensor-status__name ${statusClass}`}>
            {sensor.Sensor}
          </div>
          <div className={`sensor-status__status ${statusClass}`}>
            Setpoint: {sensor.SetPoint}
          </div>
          {sensor.Error && (
            <div className={`sensor-status__status ${statusClass}`}>
              Error: {sensor.Error}
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default SensorStatus;
