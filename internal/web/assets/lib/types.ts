interface Route {
	name: string;
	url: string;
	time: string;
	source_host: string;
}

interface Msg {
	ok: boolean;
	error?: string;
}

interface MsgRoute extends Msg {
	route: Route;
}