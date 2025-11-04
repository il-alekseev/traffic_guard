export type Device = string;

export interface DeviceControlState {
  devices: Device[] | null;
}

export interface SessionControlActions {
  fetchDevices(): Promise<Device[]>
}
