import { Button, Icon } from 'semantic-ui-react';
import { User } from '../types';
import { UserMini } from './UserBanners';

export type CommentProps = {
  author: User;
  content: string;
};

export function Comment({ author, content }: CommentProps) {
  const onReport = () => {};

  return (
    <div className='comment-container'>
      <div className='comment-user'>
        <UserMini user={author}/>
      </div>
      <div className='comment-content'>
        <div className='comment-content-text'>{content}</div>
        <div className='comment-content-actions'>
          <Button icon onClick={onReport}><Icon color='red' name="flag" /></Button>
        </div>
      </div>
    </div>
  );
}
