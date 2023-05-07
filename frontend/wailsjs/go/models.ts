export namespace main {
	
	export class ProfileMetadata {
	    name?: string;
	    about?: string;
	    picture?: string;
	    nip05?: string;
	    display_name: string;
	    lud06?: string;
	    lud16?: string;
	    banner: string;
	    website: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfileMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.about = source["about"];
	        this.picture = source["picture"];
	        this.nip05 = source["nip05"];
	        this.display_name = source["display_name"];
	        this.lud06 = source["lud06"];
	        this.lud16 = source["lud16"];
	        this.banner = source["banner"];
	        this.website = source["website"];
	    }
	}
	export class Profile {
	    pk: string;
	    following: boolean;
	    meta: ProfileMetadata;
	    npub: string;
	    relays: string[];
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pk = source["pk"];
	        this.following = source["following"];
	        this.meta = this.convertValues(source["meta"], ProfileMetadata);
	        this.npub = source["npub"];
	        this.relays = source["relays"];
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
	
	export class RelayStruct {
	    url: string;
	    read: boolean;
	    write: boolean;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RelayStruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.read = source["read"];
	        this.write = source["write"];
	        this.enabled = source["enabled"];
	    }
	}

}

export namespace nostr {
	
	export class Event {
	    id: string;
	    pubkey: string;
	    created_at: number;
	    kind: number;
	    tags: string[][];
	    content: string;
	    sig: string;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.pubkey = source["pubkey"];
	        this.created_at = source["created_at"];
	        this.kind = source["kind"];
	        this.tags = source["tags"];
	        this.content = source["content"];
	        this.sig = source["sig"];
	    }
	}

}

