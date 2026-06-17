export namespace main {
	
	export class DialogOptions {
	    defaultPath: string;
	
	    static createFrom(source: any = {}) {
	        return new DialogOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultPath = source["defaultPath"];
	    }
	}
	export class FileInfo {
	    content: string;
	    filePath: string;
	    fileName: string;
	    dirName: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.filePath = source["filePath"];
	        this.fileName = source["fileName"];
	        this.dirName = source["dirName"];
	    }
	}
	export class UserInfo {
	    reporterName: string;
	    reporterNameSig: string;
	    position: string;
	    office: string;
	    supervisorName: string;
	    supervisorPos1: string;
	    supervisorPos2: string;
	
	    static createFrom(source: any = {}) {
	        return new UserInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reporterName = source["reporterName"];
	        this.reporterNameSig = source["reporterNameSig"];
	        this.position = source["position"];
	        this.office = source["office"];
	        this.supervisorName = source["supervisorName"];
	        this.supervisorPos1 = source["supervisorPos1"];
	        this.supervisorPos2 = source["supervisorPos2"];
	    }
	}
	export class OngoingTask {
	    name: string;
	    percentComplete: number;
	    status: string;
	    nextSteps: string[];
	
	    static createFrom(source: any = {}) {
	        return new OngoingTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.percentComplete = source["percentComplete"];
	        this.status = source["status"];
	        this.nextSteps = source["nextSteps"];
	    }
	}
	export class ReportEntry {
	    date: string;
	    tasks: string[];
	    hoursWorked: number;
	    timeIn: string;
	    timeOut: string;
	
	    static createFrom(source: any = {}) {
	        return new ReportEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.tasks = source["tasks"];
	        this.hoursWorked = source["hoursWorked"];
	        this.timeIn = source["timeIn"];
	        this.timeOut = source["timeOut"];
	    }
	}
	export class ParsedData {
	    reportType: string;
	    dateStart: string;
	    dateEnd: string;
	    entries: ReportEntry[];
	    summary: string;
	    totalHours: number;
	    ongoingTasks: OngoingTask[];
	
	    static createFrom(source: any = {}) {
	        return new ParsedData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.reportType = source["reportType"];
	        this.dateStart = source["dateStart"];
	        this.dateEnd = source["dateEnd"];
	        this.entries = this.convertValues(source["entries"], ReportEntry);
	        this.summary = source["summary"];
	        this.totalHours = source["totalHours"];
	        this.ongoingTasks = this.convertValues(source["ongoingTasks"], OngoingTask);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class GenerateReportPayload {
	    parsedData: ParsedData;
	    userInfo: UserInfo;
	    employeeType: string;
	    outputPath: string;
	
	    static createFrom(source: any = {}) {
	        return new GenerateReportPayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.parsedData = this.convertValues(source["parsedData"], ParsedData);
	        this.userInfo = this.convertValues(source["userInfo"], UserInfo);
	        this.employeeType = source["employeeType"];
	        this.outputPath = source["outputPath"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	
	export class ParsePayload {
	    content: string;
	    dateStart: string;
	    dateEnd: string;
	    reportType: string;
	
	    static createFrom(source: any = {}) {
	        return new ParsePayload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.dateStart = source["dateStart"];
	        this.dateEnd = source["dateEnd"];
	        this.reportType = source["reportType"];
	    }
	}
	
	

}

