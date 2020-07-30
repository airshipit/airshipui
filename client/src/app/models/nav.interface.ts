export class NavInterface {
  displayName: string;
  disabled?: boolean;
  iconName?: string;
  route?: string;
  external?: boolean;
  children?: NavInterface[];
}
