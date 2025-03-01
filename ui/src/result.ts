export interface Result<T> {
	value: T;
	error: string;
}

export const defaultErrorToString = (e: unknown): string => {
	if (typeof e === 'string') {
		return e;
	} else if (e instanceof Error) {
		return e.message;
	}
	return 'An unknown error occurred';
}

export const of = <T>(value: T): Result<T> => {
	return { value, error: '' };
}

export const from = async <T>(
	op: () => Promise<T>,
	defaultValue: T,
	errorToString: (e: unknown) => string = defaultErrorToString,
): Promise<Result<T>> => {
	try {
		return of(await op());
	} catch (e) {
		return { value: defaultValue, error: errorToString(e) };
	}
}

export const Result = { of, from };