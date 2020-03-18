interface Route {
	name: string;
	url: string;
	time: string;
	hits: string;
}

interface Msg {
	ok: boolean;
	error?: string;
}

interface MsgRoute extends Msg {
	route: Route;
}