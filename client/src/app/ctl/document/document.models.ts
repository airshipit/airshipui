export class KustomNode {
    id: string;
    name: string;
    canLoadChildren: boolean;
    children: KustomNode[];
    isPhaseNode: boolean;
    hasError: boolean;
}
