import { useEffect, type RefObject } from 'react';


export const useResizeTextarea = (
  ref: RefObject<HTMLTextAreaElement>,
  value: string,
  maxHeight: number = 150,
  autoFocus: boolean = true,
  minHeight: number = 44,
) => {
  useEffect(() => {
    if (autoFocus && ref.current) {
      ref.current.focus();
    }
  }, [autoFocus, ref]);

  useEffect(() => {
    const textarea = ref.current;
    if (!textarea) return;

    textarea.style.height = 'auto';
    
    const scrollHeight = textarea.scrollHeight;
    const targetHeight = Math.max(minHeight, Math.min(scrollHeight, maxHeight));

    textarea.style.height = `${targetHeight}px`;

    textarea.style.overflowY = scrollHeight > maxHeight ? 'auto' : 'hidden';
    
  }, [ref, value, minHeight, maxHeight]);
};


;
