import { Button, Dropdown, Grid, GridColumn, Icon, Input, Modal, Segment } from 'semantic-ui-react';
import { ContentReport } from '../types';
import { easyDateTimeFormat } from '../utils/misc';
import { UserMini } from './UserBanners';
import { renderReportState } from './pages/Moderation';
import { useNavigate } from 'react-router-dom';
import { useMemo } from 'react';

export type ContentReportBoxProps = {
  contentReport: ContentReport | null;
  onClose: () => void;
};

type ContentRef = {
  type: string;
  id: string;
};

export function ContentReportModal(props: ContentReportBoxProps) {
  const navigate = useNavigate();

  const contentRefs = useMemo(() => {
    if (props.contentReport && props.contentReport.contentRef) {
      return parseContextRef(props.contentReport.contentRef);
    } else {
      return [];
    }
  }, [props.contentReport]);

  const visitContent = () => {
    if (props.contentReport === null) return;

    if (contentRefs.length === 1) {
      switch (contentRefs[0].type) {
        case 'playlist':
          navigate(`/playlist/${contentRefs[0].id}`);
          break;
        default:
          alert('Unknown content type');
      }
    }
  };

  return (
    <Modal
      size='large'
      onClose={() => props.onClose()}
      open={props.contentReport !== null}>
      { props.contentReport && (
        <>
          <Modal.Header>Content Report</Modal.Header><Modal.Content>
            <Grid>
              <GridColumn width={8}>
                <Segment>
                  <Grid>
                    <GridColumn width={16}>
                      <h3>Accused User</h3>
                      <UserMini user={props.contentReport.reportedUser} />
                    </GridColumn>
                    <GridColumn width={16}>
                      <h3>Reported By</h3>
                      <UserMini user={props.contentReport.reportedBy} />
                    </GridColumn>
                    <GridColumn width={16}>
                      <h3>Reported At</h3>
                      <p>{easyDateTimeFormat(props.contentReport.createdAt)}</p>
                    </GridColumn>
                    <GridColumn width={16}>
                      <h3>Report Reason</h3>
                      <p>{props.contentReport.reportReason}</p>
                    </GridColumn>
                    <GridColumn width={16}>
                      <h3>Additional Context</h3>
                      <p>{props.contentReport.context || 'None'}</p>
                    </GridColumn>
                    <GridColumn width={16}>
                      <h3>Reported Content Type</h3>
                      <p>{props.contentReport.contentRef}</p>
                    </GridColumn>
                    <GridColumn width={16}>
                      <Button onClick={visitContent} size='big'>Visit Content</Button>
                    </GridColumn>
                  </Grid>
                </Segment>
              </GridColumn>
              <GridColumn width={8}>
                <Segment>
                  <Grid>
                    <GridColumn width={16}>
                      <h3>State</h3>
                      <p>{renderReportState(props.contentReport.state, true)}</p>
                    </GridColumn>
                    <GridColumn width={16}>
                      <h3>Resolved By</h3>
                      {props.contentReport.resolvedBy.id ? <UserMini user={props.contentReport.resolvedBy} /> : 'N/A'}
                    </GridColumn>
                    <GridColumn width={16}>
                      <h3>Resolved At</h3>
                      <p>{props.contentReport.resolvedAt ? easyDateTimeFormat(props.contentReport.resolvedAt) : 'N/A'}</p>
                    </GridColumn>
                    <GridColumn width={16}>
                      <h3>Action Taken</h3>
                      <p>{props.contentReport.actionTaken ? props.contentReport.actionTaken : 'N/A'}</p>
                    </GridColumn>
                    { props.contentReport.state === 'reported' && (
                      <>
                        <GridColumn width={16}>
                          <Dropdown
                            placeholder='Select Action'
                            fluid
                            selection>
                            <Dropdown.Menu>
                              <Dropdown.Item>None</Dropdown.Item>
                              <Dropdown.Item><Icon color='yellow' name='warning sign'/> Warn</Dropdown.Item>
                              <Dropdown.Item><Icon color='red' name='ban'/> Ban</Dropdown.Item>
                            </Dropdown.Menu>
                          </Dropdown>
                        </GridColumn>
                        <GridColumn width={16}>
                          <Input placeholder='(Optional) Message to the accused' fluid/>
                        </GridColumn>
                        <GridColumn width={16}>
                          <Input placeholder='Extra context (Visible to mods only)' fluid/>
                        </GridColumn>
                        <GridColumn width={16}>
                          <Button primary>Resolve</Button>
                        </GridColumn>
                      </>
                    )}
                  </Grid>
                </Segment>
              </GridColumn>
            </Grid>
          </Modal.Content>
        </>
      )}
    </Modal>
  );
}

function parseContextRef(ref: string): ContentRef[] {
  const parts = ref.split(':');
  return parts.map((part) => {
    const [type, id] = part.split('_');
    return { type, id };
  });
}
