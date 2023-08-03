// Interface for the data to plot
export interface IDataPoint {
  Time: number;
  Value: number;
}

// Interface for the data returned from the API
export interface IData {
  Sensor: string;
  SetPoint: number;
  SetPointTs: number;
  Value: number;
  Ts: number;
  Error: null;
  TsSent: number;
}
