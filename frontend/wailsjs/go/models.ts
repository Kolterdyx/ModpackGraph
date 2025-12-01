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
	    label?: string;
	    required?: boolean;
	}
	export interface FileFilter {
	    displayName: string;
	    pattern: string;
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
	export interface Node {
	    id?: string | number;
	    name?: string;
	    icon?: string;
	    present?: boolean;
	    presentVersion?: string;
	    requiredVersion?: string;
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

export namespace keys {
	
	export interface Accelerator {
	    Key: string;
	    Modifiers: string[];
	}

}

export namespace menu {
	
	export interface MenuItem {
	    Label: string;
	    Role: number;
	    Accelerator?: keys.Accelerator;
	    Type: string;
	    Disabled: boolean;
	    Hidden: boolean;
	    Checked: boolean;
	    SubMenu?: Menu;
	}
	export interface Menu {
	    Items: MenuItem[];
	}

}

