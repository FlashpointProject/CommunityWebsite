import { Checkbox, CheckboxProps, Dropdown, DropdownProps, Grid, Segment } from 'semantic-ui-react';
import { FilterPlaylists } from '../types';
import { verifyAdult } from '../utils/adult';
import { RootState } from '../store';
import { useDispatch, useSelector } from 'react-redux';

type DropdownOption = {
  key: string;
  text: string;
  value: string;
};


const orderOptions: DropdownOption[] = [
  { key: 'name', text: 'Name', value: 'name' },
  { key: 'total_games', text: 'Total Games', value: 'total_games' },
  { key: 'created_at', text: 'Created At', value: 'created_at' },
  { key: 'updated_at', text: 'Last Updated', value: 'updated_at' }
];

const libraryOptions: DropdownOption[] = [
  { key: 'all', text: 'All', value: '' },
  { key: 'arcade', text: 'Games', value: 'arcade' },
  { key: 'theatre', text: 'Animations', value: 'theatre' }
];

type FilterBarPlaylistsProps = {
  filter: FilterPlaylists;
  onChange: (filter: FilterPlaylists) => void;
  style?: React.CSSProperties;
};

export function FilterBarPlaylists({ filter, onChange, style }: FilterBarPlaylistsProps) {
  const { adult } = useSelector((state: RootState) => state.mainState);
  const dispatch = useDispatch();

  const makeFilterDropdownHandler = (key: keyof FilterPlaylists) => {
    return (event: React.SyntheticEvent<HTMLElement, Event>, data: DropdownProps) => {
      const newFilter = { ...filter, [key]: data.value };
      onChange(newFilter);
    };
  };
  const makeCheckboxHandler = (key: keyof FilterPlaylists) => {
    return (event: React.FormEvent<HTMLInputElement>, data: CheckboxProps) => {
      const newFilter = { ...filter, [key]: data.checked };
      onChange(newFilter);
    };
  };

  return (
    <div style={style}>
      <Segment>
        <Grid columns={3}>
          <Grid.Row className='filter-headers'>
            <Grid.Column>
              <div className='filter-title'>Order By</div>
            </Grid.Column>
            <Grid.Column>
              <div className='filter-title'>Library</div>
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
                  onChange={makeFilterDropdownHandler('order')}
                  value={filter.order}
                  selection />
                <Checkbox
                  toggle
                  checked={filter.orderReverse}
                  onChange={makeCheckboxHandler('orderReverse')}
                  label='Reverse Order'/>
              </div>
            </Grid.Column>
            <Grid.Column>
              <Dropdown
                fluid
                placeholder='All'
                name='library'
                options={libraryOptions}
                onChange={makeFilterDropdownHandler('library')}
                value={filter.library}
                selection />
            </Grid.Column>
            <Grid.Column>
              <Checkbox
                toggle
                checked={filter.extreme}
                onChange={(...args) => verifyAdult(adult, dispatch, makeCheckboxHandler('extreme'), ...args)}
                label='Include Adult Content' />
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </Segment>
    </div>
  );
}
