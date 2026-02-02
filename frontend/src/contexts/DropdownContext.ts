import { createContext, useContext } from "react";

interface DropdownCompContextType {
    clickedItemId: number | null;
    shouldExpand: boolean;
    openModal: () => void;
    clickOnItem: (id: number | null) => void;
    closeModal: () => void;
    toggleModal: () => void;
    goBack: () => void;
    canGoBack: boolean;
}

export const DropdownCompContext = createContext<DropdownCompContextType | null>({
    clickedItemId: null,
    shouldExpand: false,
    openModal: () => { },
    clickOnItem: () => { },
    closeModal: () => { },
    toggleModal: () => { },
    goBack: () => { },
    canGoBack: false,
});

export const useDropdownCompContext = () => {
    const context = useContext(DropdownCompContext);
    if (!context) throw new Error("DropdownCompContext used outside of bounds");
    return context;
}