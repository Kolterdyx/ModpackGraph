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
	export interface FileFilter {
	    displayName: string;
	    pattern: string;
	}
	export interface GraphOptions {
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

