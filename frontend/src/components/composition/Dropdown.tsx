import React, { useRef, useState, forwardRef, useEffect, type ReactElement } from 'react';
import styled from "@emotion/styled";
import { DropdownCompContext, useDropdownCompContext } from '../../contexts/DropdownContext';
import useOutsideClick from '../../hooks/useOutsideClick';
import Flex from '../base/Flex';

interface ItemProps extends React.ButtonHTMLAttributes<HTMLElement> {
    children: React.ReactNode;
    isLeafItem?: boolean;
    selectedColor?: string;
    selectedBgColor?: string;
    padding?: string;
    margin?: string;
    height?: string;
    width?: string;
    color?: string;
    bgColor?: string;
    clickHandler?: (...args: any[]) => any | Promise<any>;
    borderRadius?: string;
    targetMenuIndex?: number | null;
    hoverColor?: string;
    notClickable?: boolean;
}

function Item({
    children,
    isLeafItem = false,
    selectedBgColor,
    selectedColor,
    clickHandler,
    notClickable,
    padding,
    margin,
    height,
    borderRadius,
    targetMenuIndex,
    hoverColor,
    width,
    color,
    bgColor,
    ...RestProps
}: ItemProps) {
    const { clickOnItem, clickedItemId } = useDropdownCompContext();

    return (
        <ItemContainer margin={margin}>
            {!notClickable ? (
                <ItemButton
                    padding={padding}
                    hoverColor={hoverColor}
                    height={height}
                    width={width}
                    borderRadius={borderRadius}
                    selectedColor={selectedColor}
                    selectedBgColor={selectedBgColor}
                    isSelected={clickedItemId === targetMenuIndex}
                    color={color}
                    bgColor={bgColor}
                    onClick={(e: React.MouseEvent) => {
                        e.preventDefault();
                        e.stopPropagation();
                        clickOnItem(targetMenuIndex ?? null);
                        clickHandler && clickHandler();
                    }}
                    notClickable={notClickable}
                    {...RestProps}
                >
                    {children}
                </ItemButton>
            ) : (
                <Flex justify="center" align="center" {...RestProps}>
                    {children}
                </Flex>
            )}
        </ItemContainer>
    );
}

const ItemContainer = styled.li<{ margin?: string }>`
    list-style-type: none;
    margin: ${props => props.margin};
`;

const ItemButton = styled.button<{
    padding: string;
    height?: string;
    width: string;
    borderRadius: string;
    hoverColor?: string;
    color?: string;
    bgColor?: string;
    selectedColor?: string;
    selectedBgColor?: string;
    notClickable?: boolean;
    isSelected?: boolean;
}>(props => ({
    borderRadius: props.borderRadius,
    padding: props.padding,
    height: props.height,
    outline: 'none',
    boxShadow: 'none',
    width: props.width,
    flexDirection: 'column',
    display: 'flex',
    cursor: props.notClickable ? 'default' : "pointer",
    border: 'none',
    transition: 'background-color 0.3s',
    fontFamily: 'inherit',
    color: props.isSelected ? props.selectedColor : props.color,
    backgroundColor: props.isSelected ? props.selectedBgColor : props.bgColor,
    '&:hover': {
        backgroundColor: !props.notClickable ? props.hoverColor : undefined,
    },
    '&:focus': {
        outline: 'none',
    }
}));

interface MenuProps extends React.HTMLAttributes<HTMLUListElement> {
    children: React.ReactNode | React.ReactNode[];
    _top?: string;
    plus?: number;
    bottom?: string;
    left?: string;
    right?: string;
    height?: string;
    width?: string;
    padding?: string;
    itemGap?: string;
    borderRadius?: string;
    border?: string;
    margin?: string;
    menuIndex: number;
}

function Menu({
    children,
    margin,
    _top,
    plus,
    bottom,
    left,
    right,
    height,
    width,
    padding,
    itemGap,
    borderRadius,
    border,
    menuIndex,
    ...RestProps
}: MenuProps) {
    return (
        <MenuContainer
            _top={_top}
            bottom={bottom}
            left={left}
            plus={plus}
            right={right}
            height={height}
            width={width}
            padding={padding}
            itemGap={itemGap}
            borderRadius={borderRadius}
            border={border}
            margin={margin}
            {...RestProps}
        >
            {React.Children.map(children, (child) => {
                if (React.isValidElement(child)) {
                    return React.cloneElement(child);
                }
                return null;
            })}
        </MenuContainer>
    );
}

Menu.Item = Item;

const MenuContainer = styled.ul<{
    _top?: string;
    plus?: number;
    bottom?: string;
    left?: string;
    right?: string;
    height?: string;
    width?: string;
    padding?: string;
    itemGap?: string;
    borderRadius?: string;
    border?: string;
    margin?: string;
}>(props => ({
    position: 'absolute',
    overflowX: 'hidden',
    overflowY: 'auto',
    top: props._top == undefined ? `calc( 100% + ${props.plus ? props.plus : 4}px )` : props._top,
    bottom: props.bottom,
    left: props.left,
    right: props.right,
    height: props.height,
    width: props.width,
    padding: props.padding,
    border: props.border,
    backgroundColor: 'white',
    gap: props.itemGap,
    borderRadius: props.borderRadius,
    zIndex: 999999,
    margin: props.margin || 0,
}));

interface MenuListProps {
    children: ReactElement<MenuProps> | ReactElement<MenuProps>[];
}

function MenuList({ children }: MenuListProps) {
    return (
        <MenuListContainer>
            {children}
        </MenuListContainer>
    );
}

MenuList.Menu = Menu;

const MenuListContainer = styled.div`
    margin: 0 0;
`;

interface TargetProps extends React.HTMLAttributes<HTMLDivElement> {
    children: React.ReactNode;
}

const Target = forwardRef<HTMLDivElement, TargetProps>(({ children, ...RestProps }: TargetProps, ref) => {
    const { toggleModal } = useDropdownCompContext();

    return (
        <TargetContainer ref={ref} onClick={() => toggleModal()} {...RestProps}>
            {children}
        </TargetContainer>
    );
});

const TargetContainer = styled.div``;

interface BackButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    children?: React.ReactNode;
    padding?: string;
    width?: string;
    borderRadius?: string;
    hoverColor?: string;
    color?: string;
    bgColor?: string;
}

const BackButton = ({
    children,
    padding = "0.5rem",
    width = "100%",
    borderRadius = "7px",
    hoverColor = "#f0f0f0",
    color = "black",
    bgColor = "white",
    ...RestProps
}: BackButtonProps) => {
    const { goBack, canGoBack } = useDropdownCompContext();

    if (!canGoBack) return null;

    return (
        <ItemContainer>
            <ItemButton
                padding={padding}
                width={width}
                borderRadius={borderRadius}
                hoverColor={hoverColor}
                color={color}
                bgColor={bgColor}
                onClick={(e: React.MouseEvent) => {
                    e.preventDefault();
                    e.stopPropagation();
                    goBack();
                }}
                {...RestProps}
            >
                {children || (
                    <Flex justify="flex-start" align="center" style={{ gap: "0.6rem" }}>
                        <span>←</span>
                        <span>Retour</span>
                    </Flex>
                )}
            </ItemButton>
        </ItemContainer>
    );
};

interface DropdownProps extends React.HTMLAttributes<HTMLDivElement> {
    children?: (React.ReactNode | ReactElement<MenuListProps> | undefined)[];
    withoutTarget?: boolean;
    shouldOpen?: boolean;
}

const Dropdown = ({
    children,
    withoutTarget = false,
    shouldOpen,
    ...RestProps
}: DropdownProps) => {
    const [clickedItemId, setClickedItemId] = useState<number | null>(null);
    const [shouldExpand, setShouldExpand] = useState<boolean>(false);
    const [menuHistory, setMenuHistory] = useState<number[]>([0]);
    const [Target, MenuList] = React.Children.toArray(children) as ReactElement[];
    const ref = useRef<HTMLDivElement>(null);

    const menuList = React.useMemo(() => {
        if (!React.isValidElement(MenuList)) return [];

        const menuListChildren = MenuList.props.children;
        return React.Children.toArray(menuListChildren);
    }, [MenuList]);

    const openModal = () => {
        setShouldExpand(true);
        setMenuHistory([0]);
    };

    const closeModal = () => {
        setShouldExpand(false);
        setClickedItemId(null);
        setMenuHistory([0]);
    };

    const clickOnItem = (targetIndex: number | null) => {
        if (targetIndex === null) return;

        if (targetIndex >= menuList.length || targetIndex < 0) {
            console.warn(`Menu avec l'index ${targetIndex} n'existe pas. Menus disponibles: 0-${menuList.length - 1}`);
            return;
        }

        setClickedItemId(targetIndex);

        if (targetIndex > 0) {
            setMenuHistory(prev => [...prev, targetIndex]);
            setShouldExpand(true);
        }
    };

    const goBack = () => {
        if (menuHistory.length > 1) {
            setMenuHistory(prev => prev.slice(0, -1));
            setClickedItemId(null);
        } else {
            closeModal();
        }
    };

    const toggleModal = () => {
        if (shouldExpand) {
            closeModal();
        } else {
            openModal();
        }
    };

    const currentMenuIndex = menuHistory[menuHistory.length - 1];
    const activeMenu = React.useMemo(() => {
        if (!shouldExpand && !shouldOpen) return null;

        return menuList[currentMenuIndex] || null;
    }, [shouldExpand, shouldOpen, currentMenuIndex, menuList]);

    useEffect(() => {
        if (shouldExpand !== undefined) {
            setShouldExpand(shouldExpand);
            if (!shouldExpand) {
                setMenuHistory([0]);
                setClickedItemId(null);
            }
        }
    }, [shouldExpand]);

    useOutsideClick(ref, closeModal);

    return (
        <DropdownCompContext.Provider value={{
            clickedItemId,
            shouldExpand,
            openModal,
            closeModal,
            clickOnItem,
            toggleModal,
            goBack,
            canGoBack: menuHistory.length > 1,
        }}>
            <DropdownCompContainer
                ref={ref}
                {...RestProps}
            >
                {!withoutTarget && Target}
                {activeMenu}
            </DropdownCompContainer>
        </DropdownCompContext.Provider>
    );
};

const DropdownCompContainer = styled.div({
    position: "relative",
});

Dropdown.Target = Target;
Dropdown.MenuList = MenuList;
Dropdown.BackButton = BackButton;

export default Dropdown;
