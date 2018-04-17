interface Route {
	name: string;
	url: string;
	uid: string;
    created_at: string;
    modified_at: string;
    deleted_at: string;
	generated: boolean;
	modified_count: number;
}

interface Msg {
	ok: boolean;
	error?: string;
}

interface MsgRoute extends Msg {
	route: Route;
}