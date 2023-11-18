import { NewsPost } from '../types';
import { FlashpointMarkdown } from './FlashpointMarkdown';
import { SegmentWithHeader } from './SegmentWithHeader';
import { UserMini } from './UserBanners';

type NewsPostBoxProps = {
  post: NewsPost;
};

/**
 * A segment with a fancy header to display like a news post
 * @param title Title of the post
 * @param user Author of the post
 * @param content Content of the post
 */
export const NewsPostBox = ({ post }: NewsPostBoxProps) => {
  const header = (
    <>
      <div className="left-flex">
        {post.title}
      </div>
      <div className="right-flex">
        <UserMini user={post.author} />
      </div>
    </>
  );

  return (
    <SegmentWithHeader header={header}>
      <FlashpointMarkdown>{post.content}</FlashpointMarkdown>
    </SegmentWithHeader>
  );
};

