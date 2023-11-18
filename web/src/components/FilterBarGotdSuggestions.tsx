import { Checkbox, Dropdown, Grid, Segment } from 'semantic-ui-react';
import { FilterGotdSuggestions } from '../types';
import { makeCheckboxHandler, makeDropdownHandler } from '../utils/filterBarCommon';

type DropdownOption = {
  key: string;
  text: string;
  value: string;
};

const orderOptions: DropdownOption[] = [
  { key: 'created_at', text: 'Suggested At', value: 'created_at' },
  { key: 'suggested_date', text: 'Suggested Date', value: 'suggested_date' },
];

type FilterBarGotdSuggestionsProps = {
  filter: FilterGotdSuggestions;
  onChange: (filter: FilterGotdSuggestions) => void;
  style?: React.CSSProperties;
};

export function FilterBarGotdSuggestions({ filter, onChange, style }: FilterBarGotdSuggestionsProps) {
  return (
    <div style={style}>
      <Segment>
        <Grid columns={2}>
          <Grid.Row className='filter-headers'>
            <Grid.Column>
              <div className='filter-title'>Order By</div>
            </Grid.Column>
          </Grid.Row>
          <Grid.Row>
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
