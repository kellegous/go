import { createRef, useEffect, useState } from "react";
import css from "./Copy.module.scss";

export interface CopyProps {
  text: string;
}

export const Copy = ({ text }: CopyProps) => {
  const [active, setActive] = useState(false);
  const ref = createRef<HTMLDivElement>();

  useEffect(() => {
    const { current } = ref;
    if (active && current) {
      const animationDidEnd = () => {
        setActive(false);
      };
      current.addEventListener("animationend", animationDidEnd);
      return () => {
        current.removeEventListener("animationend", animationDidEnd);
      };
    }
  }, [active, setActive, ref]);

  const copyWasClicked = () => {
    navigator.clipboard.writeText(text);
    setActive(true);
  };

  const feedbackClasses = active
    ? `${css.feedback} ${css.active} material-symbols-outlined`
    : `${css.feedback} material-symbols-outlined`;

  return (
    <div className={css.button}>
      <div className={feedbackClasses} ref={ref}>
        task_alt
      </div>
      <div
        className={`${css.icon} material-symbols-outlined`}
        onClick={copyWasClicked}
      >
        content_copy
      </div>
    </div>
  );
};
