import css from "./Controls.module.scss";

export interface ControlsProps {
  name: string;
}

export const Controls = ({ name }: ControlsProps) => {
  return (
    <a href={`/edit/${name}`} className={css.controls}>
      <span className="material-symbols-outlined">edit</span>
    </a>
  );
};
