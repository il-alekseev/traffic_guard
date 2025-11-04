import { Categories } from "./categories";


export interface CategoriesControlState {
  categories: Categories[] | null;
}

export interface DeviceontrolActions {
  fetchCategories(): Promise<Categories[]>
}
