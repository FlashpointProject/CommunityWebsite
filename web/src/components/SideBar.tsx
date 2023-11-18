import { SegmentWithHeader } from './SegmentWithHeader';

type SideBarProps = {
  title: string;
  icon?: string;
} & React.PropsWithChildren;

/**
 * An alias to SegmentWithHeader without fill
 * @param title Title of the sidebar
 * @param icon (Optional) Icon to prefix the title with
 */
export function SideBar(props: SideBarProps) {
  const header = (
    <>
      {props.icon && <img src={props.icon} alt={props.title} />}
      <h3>{props.title}</h3>
    </>
  );

  return (
    <SegmentWithHeader header={header}>
      {props.children}
    </SegmentWithHeader>
  );
}
