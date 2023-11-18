import { Checkbox, Dropdown, Grid, Segment } from 'semantic-ui-react';
import { FilterContentReports } from '../types';
import { makeCheckboxHandler, makeDropdownHandler } from '../utils/filterBarCommon';

type DropdownOption = {
  key: string;
  text: string;
  value: string;
};

const orderOptions: DropdownOption[] = [
  { key: 'content_type', text: 'Content Type', value: 'content_type' },
  { key: 'aggregate_report_score', text: 'Aggregate Report Score', value: 'aggregate_report_score' },
  { key: 'created_at', text: 'Created At', value: 'created_at' },
  { key: 'updated_at', text: 'Last Updated', value: 'updated_at' },
  { key: 'resolved_at', text: 'Resolved At', value: 'resolved_at' }
];

const reportStateOptions = [
  { key: 'all', text: 'All', value: '' },
  { key: 'reported', text: 'Reported', value: 'reported', icon: 'flag'},
  { key: 'resolved', text: 'Resolved', value: 'resolved', icon: 'checkmark' },
];

type FilterBarContentReportsProps = {
  filter: FilterContentReports;
  onChange: (filter: FilterContentReports) => void;
  style?: React.CSSProperties;
};

export function FilterBarContentReports({ filter, onChange, style }: FilterBarContentReportsProps) {
  return (
    <div style={style}>
      <Segment>
        <Grid columns={2}>
          <Grid.Row className='filter-headers'>
            <Grid.Column>
              <div className='filter-title'>Report State</div>
            </Grid.Column>
            <Grid.Column>
              <div className='filter-title'>Order By</div>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
            <Grid.Column>
              <Dropdown
                placeholder='All'
                options={reportStateOptions}
                value={filter.reportState}
                onChange={makeDropdownHandler('reportState', filter, onChange)}
                selection>
              </Dropdown>
            </Grid.Column>
            <Grid.Column>
              <div className='filter-playlists-combined-column'>
                <Dropdown
                  fluid
                  placeholder='Order By'
                  name='orderBy'
                  options={orderOptions}
                  onChange={makeDropdownHandler('order', filter, onChange)}
                  value={filter.order}
                  selection />
                <Checkbox
                  toggle
                  checked={filter.orderReverse}
                  onChange={makeCheckboxHandler('orderReverse', filter, onChange)}
                  label='Reverse Order'/>
              </div>
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </Segment>
    </div>
  );
}
