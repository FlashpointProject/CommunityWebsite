import { Header, Segment } from 'semantic-ui-react';

type SegmentWithHeaderProps = {
  header: React.ReactNode;
  fill?: boolean;
} & React.PropsWithChildren;

/**
 * A segment with a distinct header
 * @param header Header of the segment
 * @param children Content of the segment
 * @param fill Whether the segment should fill the entire width and height of the page
 */
export function SegmentWithHeader(props: SegmentWithHeaderProps) {
  return (
    <Segment className={`segment-with-header ${props.fill ? 'segment-fill' : ''}`} padded={false}>
      <Header as='h3' className="segment-with-header__header block-header block-header-color">
        {props.header}
      </Header>
      <div className="segment-with-header__content">
        {props.children}
      </div>
    </Segment>
  );
}
