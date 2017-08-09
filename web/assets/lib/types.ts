interface Route {
	name: string;
	url: string;
	time: string;
}

interface Msg {
	ok: boolean;
	error?: string;
}

interface MsgRoute extends Msg {
	route: Route;
}