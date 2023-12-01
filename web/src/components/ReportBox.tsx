import { useCallback, useState } from 'react';
import { Button, CheckboxProps, Form, Modal, Radio, TextArea } from 'semantic-ui-react';

export type ReportModalProps = {
  contentRef: string;
  contentName: string;
  trigger: React.ReactNode;
};

export type ReportButtonProps = Omit<ReportModalProps, 'trigger'> & {
  withText?: boolean;
};

export function ReportModal({ trigger, contentRef, contentName }: ReportModalProps) {
  const [validReport, setValidReport] = useState(false);
  const [reportReason, setReportReason] = useState('');
  const [reportContext, setReportContext] = useState('');
  const [open, setOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const onRadioChange = (event: React.FormEvent<HTMLInputElement>, data: CheckboxProps) => {
    setReportReason(data.value.toString());
    if (data.value.toString() !== '') {
      setValidReport(true);
    } else {
      setValidReport(false);
    }
  };

  const onSubmit = useCallback(() => {
    if (submitting) {
      return;
    }
    if (validReport) {
      setSubmitting(true);
      fetch('/api/reports', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          'content_ref': contentRef,
          'reason': reportReason,
          'context': reportContext,
        }),
      }).then((response) => {
        if (response.status === 200) {
          alert('Report submitted successfully.');
          setOpen(false);
        } else {
          alert('Failed to submit report - ' + response.statusText);
        }
      })
      .catch(() => {
        alert('Failed to submit report. Server connection failed.');
      })
      .finally(() => {
        setSubmitting(false);
      });
    }
  }, [validReport, reportReason, reportContext, submitting]);

  return (
    <Modal
      open={open}
      onClose={() => setOpen(false)}
      onOpen={() => setOpen(true)}
      trigger={trigger}>
      <Modal.Header>Report Content</Modal.Header>
      <Modal.Content>
        <p>Reporting: {contentName}</p>
        <Form>
          <label>Reason for report</label>
          <Form.Field>
            <Radio
              label='Spam'
              name='report_reason'
              value='Spam'
              checked={reportReason === 'Spam'}
              onChange={onRadioChange}
            />
          </Form.Field>
          <Form.Field>
            <Radio
              label='Harassment / Hate Speech'
              name='report_reason'
              value='Harassment / Hate Speech'
              checked={reportReason === 'Harassment / Hate Speech'}
              onChange={onRadioChange}
            />
          </Form.Field>
          <Form.Field>
            <Radio
              label='Inappropriate Content'
              name='report_reason'
              value='Inappropriate Content'
              checked={reportReason === 'Inappropriate Content'}
              onChange={onRadioChange}
            />
          </Form.Field>
          <Form.Field>
            <Radio
              label='Illegal Content'
              name='report_reason'
              value='Illegal Content'
              checked={reportReason === 'Illegal Content'}
              onChange={onRadioChange}
            />
          </Form.Field>
          <Form.Field>
            <Radio
              label='Other'
              name='report_reason'
              value='Other'
              checked={reportReason === 'Other'}
              onChange={onRadioChange}
            />
          </Form.Field>
          <Form.Field>
            <label>{'Additional Context (optional)'}</label>
            <TextArea
              value={reportContext}
              onChange={(event, data) => setReportContext(data.value.toString())}
              placeholder='Additional Context' />
          </Form.Field>
        </Form>
      </Modal.Content>
      <Modal.Actions>
        <Button disabled={!validReport && !submitting} onClick={onSubmit}>Report</Button>
        <Button disabled={!submitting} onClick={() => setOpen(false)}>Cancel</Button>
      </Modal.Actions>
    </Modal>
  );
}

export function ReportButton({ withText, ...args }: ReportButtonProps) {
  const button = withText ? (
    <Button icon='flag'>Report</Button>
  ) : (
    <Button icon='flag'/>
  );

  return (
    <ReportModal
      {...args}
      trigger={button}
    />
  );
}
