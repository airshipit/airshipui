export class KustomNode {
    id: string;
    phaseid: { name: string, namespace: string};
    name: string;
    canLoadChildren: boolean;
    children: KustomNode[];
    isPhaseNode: boolean;
    hasError: boolean;
}
