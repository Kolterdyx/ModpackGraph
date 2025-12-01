export namespace app {
	
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
	}
	export interface GraphGenerationOptions {
	    path?: string;
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

