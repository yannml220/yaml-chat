import { type Dispatch, useEffect, type SetStateAction, type RefObject } from 'react'

const useOutsideClick = <T extends HTMLElement = HTMLElement>(
    ref: RefObject<T>,
    callback: Dispatch<SetStateAction<unknown>> | (() => void),
    payload?: unknown
) => {
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (ref.current && !ref.current.contains(event.target as Node)) {
                if (typeof callback === 'function') {
                    (callback as Function)(payload);
                }
            }
        }

        document.addEventListener('click', handleClickOutside);

        return () => {
            document.removeEventListener('click', handleClickOutside);
        }
    }, [ref, callback, payload]);
}

export default useOutsideClick;
