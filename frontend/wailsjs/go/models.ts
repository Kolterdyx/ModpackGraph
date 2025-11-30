export namespace app {
	
	export enum Layout {
	    circo = "circo",
	    dot = "dot",
	    fdp = "fdp",
	    neato = "neato",
	    nop = "nop",
	    nop1 = "nop1",
	    nop2 = "nop2",
	    osage = "osage",
	    patchwork = "patchwork",
	    sfdp = "sfdp",
	    twopi = "twopi",
	}
	export interface Edge {
	    source: string;
	    target: string;
	}
	export interface FileFilter {
	    displayName: string;
	    pattern: string;
	}
	export interface Node {
	    id: string;
	    color?: string;
	    name?: string;
	    val?: number;
	    icon?: string;
	}
	export interface Graph {
	    nodes: Node[];
	    links: Edge[];
	    layout?: Layout;
	}
	export interface GraphGenerationOptions {
	    path?: string;
	    layout?: Layout;
	}
	
	export interface OpenDialogOptions {
	    title?: string;
	    defaultDirectory?: string;
	    defaultFilename?: string;
	    filters?: FileFilter[];
	    showHiddenFiles?: boolean;
	    canCreateDirectories?: boolean;
	    resolvesAliases?: boolean;
	    treatPackagesAsDirectories?: boolean;
	}

}

