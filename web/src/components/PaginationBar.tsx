import { Icon, Pagination, PaginationProps, Statistic } from 'semantic-ui-react';

type PagedQuery = {
  page: number;
  pageSize: number;
};

type PaginationBarProps = {
  query: PagedQuery;
  totalResults: number;
  onPageChange: (event: React.MouseEvent<HTMLAnchorElement, MouseEvent>, data: PaginationProps) => void;
};

export function PaginationBar({ query, totalResults, onPageChange }: PaginationBarProps) {
  return (
    <div className='box-playlist-bottom-row'>
      <Pagination
        defaultActivePage={query.page}
        ellipsisItem={{ content: <Icon name='ellipsis horizontal' />, icon: true }}
        firstItem={{ content: <Icon name='angle double left' />, icon: true }}
        lastItem={{ content: <Icon name='angle double right' />, icon: true }}
        prevItem={{ content: <Icon name='angle left' />, icon: true }}
        nextItem={{ content: <Icon name='angle right' />, icon: true }}
        totalPages={Math.ceil(totalResults / query.pageSize)}
        onPageChange={onPageChange}
      />
      <Statistic size='mini'>
        <Statistic.Label>Results</Statistic.Label>
        <Statistic.Value>{totalResults}</Statistic.Value>
      </Statistic>
    </div>
  );
}
