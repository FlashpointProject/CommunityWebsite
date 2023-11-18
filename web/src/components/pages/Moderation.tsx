import React, { useCallback, useEffect, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Button, Icon, Loader, PaginationProps, Table } from 'semantic-ui-react';
import { forceContentReportsLoad, setContentReportsQuery, setContentReportsResultsPage, setOpenReport } from '../../redux/moderationSlice';
import { RootState } from '../../store';
import { ContentReport, FilterContentReports } from '../../types';
import { easyDateTimeFormat } from '../../utils/misc';
import { ContentReportModal } from '../ContentReport';
import { FilterBarContentReports } from '../FilterBarContentReports';
import { PaginationBar } from '../PaginationBar';
import { UserMini } from '../UserBanners';

export function ContentReportsPage() {
  const { query, results, totalResults, searching, openReport } = useSelector((state: RootState) => state.moderationState);
  const dispatch = useDispatch();

  const onQueryChange = useCallback((query: FilterContentReports) => {
    dispatch(setContentReportsQuery(query));
  }, []);

  const onPageChange = useCallback((event: React.MouseEvent<HTMLAnchorElement, MouseEvent>, data: PaginationProps) => {
    dispatch(setContentReportsResultsPage(data.activePage as number));
  }, [query]);

  const headerStyle = results.length > 0 ? {
    marginBottom: '0',
  } : {
    marginBottom: '1rem',
  };

  const totalPages = useMemo(() => {
    return Math.ceil(totalResults / query.pageSize);
  }, [totalResults]);

  useEffect(() => {
    dispatch(forceContentReportsLoad());
  }, []);

  const mapReportTableRow = (report: ContentReport) => {
    return (
      <Table.Row>
        <Table.Cell collapsing>
          <Button icon='linkify' onClick={() => dispatch(setOpenReport(report))} />
        </Table.Cell>
        <Table.Cell negative={report.state === 'reported'} positive={report.state === 'resolved'} collapsing>
          {renderReportState(report.state)}
        </Table.Cell>
        <Table.Cell>
          <UserMini user={report.reportedBy} />
        </Table.Cell>
        <Table.Cell>
          <UserMini user={report.reportedUser} />
        </Table.Cell>
        <Table.Cell>
          {report.resolvedBy.id ? <UserMini user={report.resolvedBy}/> : 'N/A'}
        </Table.Cell>
        <Table.Cell>
          {report.actionTaken || 'N/A'}
        </Table.Cell>
        <Table.Cell>
          {report.resolvedAt ? easyDateTimeFormat(new Date(report.resolvedAt)) : 'N/A'}
        </Table.Cell>
        <Table.Cell>
          {report.createdAt ? easyDateTimeFormat(new Date(report.createdAt)) : 'N/A'}
        </Table.Cell>
      </Table.Row>
    );
  };

  return (
    <div>
      <h3>Content Reports</h3>
      <ContentReportModal onClose={() => dispatch(setOpenReport(null))} contentReport={openReport} />
      <FilterBarContentReports filter={query} onChange={onQueryChange} style={headerStyle}/>
      { totalPages <= 0 ? searching ? (
        <Loader active inline='centered'>Searching</Loader>
      ) : (
        <h2>No Content Reports Found</h2>
      ) : (
        <>
          <Table celled padded>
            <Table.Header>
              <Table.Row>
                <Table.HeaderCell>Details</Table.HeaderCell>
                <Table.HeaderCell>State</Table.HeaderCell>
                <Table.HeaderCell>Reported By</Table.HeaderCell>
                <Table.HeaderCell>Reported User</Table.HeaderCell>
                <Table.HeaderCell>Resolved By</Table.HeaderCell>
                <Table.HeaderCell>Action Taken</Table.HeaderCell>
                <Table.HeaderCell>Resolved At</Table.HeaderCell>
                <Table.HeaderCell>Reported At</Table.HeaderCell>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {results.map(mapReportTableRow)}
            </Table.Body>
          </Table>
          <div className='box-playlist-bottom-row'>
            <PaginationBar query={query} totalResults={totalResults} onPageChange={onPageChange} />
          </div>
        </>
      )}
    </div>
  );
}

export function renderReportState(state: string, withName?: boolean) {
  switch (state) {
    case 'reported':
      return <>
        <Icon color='red' name='flag' />
        {withName ? 'Reported' : ''}
      </>;
    case 'resolved':
      return <>
        <Icon color='green' name='check' />
        {withName ? 'Resolved' : ''}
      </>;
    default:
      return <>
        <Icon color='grey' name='question' />
        {withName ? 'Unknown' : ''}
      </>;
  }
}
