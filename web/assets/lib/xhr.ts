namespace xhr {
	export class Req {
		private doneFns = [];

		private errorFns = [];

		public constructor(private xhr: XMLHttpRequest) {
			xhr.onload = () => {
				var text = xhr.responseText,
					status = xhr.status;
				this.doneFns.forEach((fn) => {
					fn(text, status);
				});
			};

			xhr.onerror = () => {
				this.errorFns.forEach((fn) => fn());
			};
		}

		public onDone(fn: (data: string, status: number) => void) {
			this.doneFns.push(fn);
			return this;
		}

		public onError(fn: () => void) {
			this.errorFns.push(fn);
			return this;
		}

		public withHeader(k: string, v: string) {
			this.xhr.setRequestHeader(k, v);
			return this;
		}

		public sendJSON(data: any) {
			this.withHeader('Content-Type', 'application/json;charset=utf8');
			this.xhr.send(JSON.stringify(data));
			return this;
		}

		public send(data?: string) {
			this.xhr.send(data);
			return this;
		}
	}

	export var create = (method: string, url: string) => {
		var xhr = new XMLHttpRequest();
		xhr.open(method, url, true);

		var req = new Req(xhr);
		return req.withHeader('X-Requested-With', 'XMLHttpRequest');
	};

	export var get = (url: string) => {
		return create('GET', url);
	}

	export var post = (url: string) => {
		return create('POST', url);
	}
}