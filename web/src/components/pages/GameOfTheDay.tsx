import { useCallback, useEffect, useMemo, useState } from 'react';
import Calendar from 'react-calendar';
import { useDispatch, useSelector } from 'react-redux';
import { Divider, Form, Loader, PaginationProps, Tab, Table } from 'semantic-ui-react';
import { forceGotdSuggestionsLoad, setGotdSuggestionsQuery, setGotdSuggestionsResultsPage } from '../../redux/gotdSlice';
import { RootState } from '../../store';
import { Game, GotdSuggestion } from '../../types';
import { mapRawGame } from '../../utils/mappers';
import { easyDateFormat } from '../../utils/misc';
import { throwOn404 } from '../../utils/permissions';
import { FilterBarGotdSuggestions } from '../FilterBarGotdSuggestions';
import { GameList } from '../GameList';
import { PaginationBar } from '../PaginationBar';

export function GotdPage() {
  const [suggestionGameId, setSuggestionGameId] = useState('');
  const [suggestionDescription, setSuggestionDescription] = useState('');
  const [suggestionAnonymous, setSuggestionAnonymous] = useState(false);
  const [gamePreview, setGamePreview] = useState<Game | null>(null);
  const [loadingGamePreview, setLoadingGamePreview] = useState(false);

  useEffect(() => {
    // Check if the game ID is valid
    if (suggestionGameId && suggestionGameId.length === 36) {
      setLoadingGamePreview(true);
      fetch(`/api/game/${suggestionGameId}`)
      .then(throwOn404)
      .then(res => res.json())
      .then((data) => {
        setGamePreview(mapRawGame(data));
      })
      .catch(console.error)
      .finally(() => {
        setLoadingGamePreview(false);
      });
    } else {
      setGamePreview(null);
    }
  }, [suggestionGameId]);

  const panes = [{
    menuItem: 'Submit',
    render: () => <Tab.Pane>
      <GotdSubmit
        loadingGamePreview={loadingGamePreview}
        gamePreview={gamePreview}
        suggestionGameId={suggestionGameId}
        suggestionDescription={suggestionDescription}
        suggestionAnonymous={suggestionAnonymous}
        setSuggestionGameId={setSuggestionGameId}
        setSuggestionDescription={setSuggestionDescription}
        setSuggestionAnonymous={setSuggestionAnonymous}
      />
    </Tab.Pane>,
  }, {
    menuItem: 'Suggestions',
    render: () => <Tab.Pane className='undo-segment'>
      <GotdSuggestions />
    </Tab.Pane>,
  }, {
    menuItem: 'Browse',
    render: () => <Tab.Pane>Browse</Tab.Pane>,
  }];

  return (
    <>
      <h3>Game of the Day</h3>
      <Tab menu={{ pointing: true }} panes={panes} />
    </>
  );
}

type GotdSubmitProps = {
  loadingGamePreview: boolean;
  gamePreview: Game | null;
  suggestionGameId: string;
  suggestionDescription: string;
  suggestionAnonymous: boolean;
  setSuggestionGameId: (gameId: string) => void;
  setSuggestionDescription: (description: string) => void;
  setSuggestionAnonymous: (anonymous: boolean) => void;
};

function GotdSubmit({ loadingGamePreview, gamePreview, suggestionGameId, suggestionDescription, suggestionAnonymous, setSuggestionGameId, setSuggestionDescription, setSuggestionAnonymous }: GotdSubmitProps) {
  return (
    <>
      <h4>Submit a suggestion for Game of the Day</h4>
      <div className='two-column-responsive-grid'>
        <div className='two-column-responsive-grid-left'>
          <Form>
            <Form.Input label='Game ID' value={suggestionGameId} onChange={(event, data) => {
              setSuggestionGameId(data.value);
            }} placeholder="Game ID" />
            <Form.TextArea label='Suggestion Description' value={suggestionDescription} onChange={(event, data) => {
              setSuggestionDescription(data.value.toString());
            }} placeholder="Suggestion Description" />
            <Form.Checkbox label='Suggest Anonymously' checked={suggestionAnonymous} onChange={(event, data) => {
              setSuggestionAnonymous(data.checked);
            }}/>
            <label>* Site Adminstrators can always see which suggestions are yours</label>
          </Form>
        </div>
        <div className='two-column-responsive-grid-right'>
          <label className='bold-label'>Suggested Date (Optional)</label>
          <Calendar
            minDate={new Date()} />
        </div>
      </div>
      <Divider />
      { loadingGamePreview ? (
        <Loader active inline>Loading Game Preview</Loader>
      ) : gamePreview ? (
        <GameList games={[gamePreview]} />
      ) : (
        <h3>No Game Found</h3>
      )}
    </>
  );
}

type GotdSuggestionsProps = object;

function GotdSuggestions(props: GotdSuggestionsProps) {
  const { results, query, totalResults, searching } = useSelector((state: RootState) => state.gotdState);
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(forceGotdSuggestionsLoad());
  }, []);

  const totalPages = useMemo(() => {
    return Math.ceil(totalResults / query.pageSize);
  }, [totalResults]);

  const onPageChange = useCallback((event: React.MouseEvent<HTMLAnchorElement, MouseEvent>, data: PaginationProps) => {
    dispatch(setGotdSuggestionsResultsPage(data.activePage as number));
  }, [query]);

  return (
    <>
      <h4>Game of the Day Suggestions</h4>
      <FilterBarGotdSuggestions filter={query} onChange={query => dispatch(setGotdSuggestionsQuery(query))}/>
      { totalPages <= 0 ?
        searching ? (
          <Loader active inline='centered'>Searching</Loader>
        ) : (
          <h3>No Suggestions Found</h3>
        ) : (
          <>
            <Table celled padded>
              <Table.Header>
                <Table.Row>
                  <Table.HeaderCell>Game</Table.HeaderCell>
                  <Table.HeaderCell>Submitted By</Table.HeaderCell>
                  <Table.HeaderCell>Description</Table.HeaderCell>
                  <Table.HeaderCell>Submitted At</Table.HeaderCell>
                  <Table.HeaderCell>Suggested Date</Table.HeaderCell>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {results.map(mapGotdSuggestionTableRow)}
              </Table.Body>
            </Table>
            <PaginationBar query={query} totalResults={totalResults} onPageChange={onPageChange} />
          </>
        )}
    </>
  );
}

function mapGotdSuggestionTableRow(data: GotdSuggestion) {
  return (
    <Table.Row>
      <Table.Cell>{data.game.title}</Table.Cell>
      <Table.Cell>{data.author}</Table.Cell>
      <Table.Cell>{data.description}</Table.Cell>
      <Table.Cell>{data.createdAt}</Table.Cell>
      <Table.Cell>{data.suggestedDate ? easyDateFormat(data.suggestedDate) : 'N/A'}</Table.Cell>
    </Table.Row>
  );
}
