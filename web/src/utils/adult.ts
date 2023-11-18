import { AnyAction, Dispatch } from '@reduxjs/toolkit';
import { openAgeVerification } from '../redux/mainSlice';

export async function verifyAdult(adult: boolean, dispatch: Dispatch<AnyAction>, callback: (...args: any[]) => void, ...args: any[]) {
  if (adult) {
    callback(...args);
  } else {
    try {
      await new Promise<void>((resolve, reject) => {
        const listener = (response: boolean) => {
          if (response) {
            resolve();
          } else {
            reject();
          }
        };
        dispatch(openAgeVerification(listener));
      });
      callback(...args);
    } catch {
      // do nothing
    }
  }
}
