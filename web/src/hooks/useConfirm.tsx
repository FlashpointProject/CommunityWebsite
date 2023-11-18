import { useCallback, useState } from 'react';
import { Button, Modal } from 'semantic-ui-react';

/** Opens a confirm dialog before executing the callback
 * @returns `modal` - The modal to display. Must be placed in the render tree.
 * @returns `confirm` - The function to call to open the modal. Pass in the title, message, and callback to execute.
 */
export function useConfirm() {
  const [confirmModalOpen, setConfirmModalOpen] = useState(false);
  const [confirmModalMessage, setConfirmModalMessage] = useState('');
  const [confirmModalTitle, setConfirmModalTitle] = useState('');
  const [confirmModalOnConfirm, setConfirmModalOnConfirm] = useState(() => () => {});

  const confirm = useCallback((title: string, message: string, onConfirm: () => void) => {
    return () => {
      setConfirmModalOpen(true);
      setConfirmModalMessage(message);
      setConfirmModalTitle(title);
      setConfirmModalOnConfirm(() => onConfirm);
    };
  }, []);

  const closeConfirmModal = useCallback(() => {
    setConfirmModalOpen(false);
    setConfirmModalMessage('');
    setConfirmModalTitle('');
    setConfirmModalOnConfirm(() => () => {});
  }, []);

  const modal = (
    <Modal
      open={confirmModalOpen}>
      <Modal.Header>{confirmModalTitle}</Modal.Header>
      <Modal.Content>
        <p>{confirmModalMessage}</p>
      </Modal.Content>
      <Modal.Actions>
        <Button onClick={closeConfirmModal} negative>Cancel</Button>
        <Button onClick={confirmModalOnConfirm} positive>Confirm</Button>
      </Modal.Actions>
    </Modal>
  );

  return {
    confirm,
    modal
  };
}
