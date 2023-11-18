import { useMemo } from 'react';
import { CheckboxProps, DropdownItemProps, DropdownProps } from 'semantic-ui-react';

type DropdownValue = boolean | number | string;
export function makeDropdownHandler<T>(
  key: { [K in keyof T]: T[K] extends DropdownValue ? K : never }[keyof T],
  filter: T,
  onChange: (newFilter: T) => void,
) {
  return useMemo(() => {
    return (event: React.SyntheticEvent<HTMLElement, Event>, data: DropdownProps | DropdownItemProps) => {
      const newFilter = { ...filter, [key]: data.value };
      onChange(newFilter);
    };
  }, [filter, key, onChange]);
}

export function makeMultiDropdownHandler<T>(
  key: { [K in keyof T]: T[K] extends DropdownValue[] ? K : never }[keyof T],
  filter: T,
  onChange: (newFilter: T) => void,
) {
  return useMemo(() => {
    return (event: React.SyntheticEvent<HTMLElement, Event>, data: DropdownItemProps) => {
      const newMultiSelect = [...filter[key] as DropdownValue[]];
      const existingIdx = newMultiSelect.indexOf(data.value);
      if (existingIdx === -1) {
        newMultiSelect.push(data.value);
      } else {
        newMultiSelect.splice(existingIdx, 1);
      }
      const newFilter = { ...filter, [key]: newMultiSelect };
      onChange(newFilter);
    };
  }, [filter, key, onChange]);
}

export function makeCheckboxHandler<T>(
  key: { [K in keyof T]: T[K] extends boolean ? K : never }[keyof T],
  filter: T,
  onChange: (newFilter: T) => void,
) {
  return useMemo(() => {
    return (event: React.SyntheticEvent<HTMLElement, Event>, data: CheckboxProps) => {
      const newFilter = { ...filter, [key]: data.checked };
      onChange(newFilter);
    };
  }, [filter, key, onChange]);
}
