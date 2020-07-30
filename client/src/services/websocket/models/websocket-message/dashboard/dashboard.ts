import {Executable} from './executable/executable';

export class Dashboard {
  name: string;
  baseURL: string;
  path: string;
  isProxied: boolean;
  executable: Executable;
}
