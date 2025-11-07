export type Device = string;

export interface DevicesState {
  devices: Device[] | null;
}

export interface DevicesActions {
  fetchDevices(): Promise<Device[]>
}
