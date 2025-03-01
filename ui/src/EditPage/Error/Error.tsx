export interface ErrorProps {
	message: string;
}

export const Error = ({ message }: ErrorProps) => {
	return (
		<div>{message}</div>
	)
}