export namespace videoManager {
	
	export class VideoFileInfo {
	    name: string;
	    id: string;
	    size: number;
	    // Go type: time
	    dateCreated: any;
	    path: string;
	    convertedPath: string;
	    converted: boolean;
	    error: string;
	    convertedSize: number;
	    width: number;
	    height: number;
	    duration: number;
	    bitrate: number;
	    framerate: string;
	
	    static createFrom(source: any = {}) {
	        return new VideoFileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.id = source["id"];
	        this.size = source["size"];
	        this.dateCreated = this.convertValues(source["dateCreated"], null);
	        this.path = source["path"];
	        this.convertedPath = source["convertedPath"];
	        this.converted = source["converted"];
	        this.error = source["error"];
	        this.convertedSize = source["convertedSize"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.duration = source["duration"];
	        this.bitrate = source["bitrate"];
	        this.framerate = source["framerate"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

