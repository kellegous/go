import css from './Drawer.module.scss';
import { DrawerStyle } from './DrawerStyle';

export interface DrawerProps {
	visible: boolean;
	style: DrawerStyle;
	children: React.ReactNode;
}

export const Drawer = ({ children, style, visible }: DrawerProps) => {
	return (
		<div className={classNameFrom(style, visible)}>
			{children}
		</div>
	);
};

const classNameFrom = (
	style: DrawerStyle,
	visible: boolean,
): string => {
	const names = [css.drawer];
	if (style === DrawerStyle.Error) {
		names.push(css.fuck);
	}
	if (visible) {
		names.push(css.visible);
	}
	return names.join(' ');
}