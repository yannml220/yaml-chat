import React, { forwardRef } from 'react'
import styled from '@emotion/styled'

interface FlexProps extends React.HTMLAttributes<HTMLDivElement> {
    height?: string,
    width?: string,
    maxHeight?: string,
    minHeight?: string,
    maxWidth?: string,
    minWidth?: string,
    padding?: string,
    margin?: string,
    color?: string,
    bgColor?: string,
    border?: string,
    radius?: string,
    shadow?: string,
    children: React.ReactNode;
    direction?: 'row' | 'column';
    justify?: string;
    align?: string;
    colGap?: string;
    rowGap?: string;
    flex?: number;
    shrink?: number;
    basis?: string;
    grow?: number;
    wrap?: 'nowrap' | 'wrap' | 'wrap-reverse';

    overflowX?: 'visible' | 'hidden' | 'clip' | 'scroll' | 'auto';
    overflowY?: 'visible' | 'hidden' | 'clip' | 'scroll' | 'auto';
    isBorderBox?: boolean;
    pointer?: string;
}




const Flex = forwardRef<HTMLDivElement, FlexProps>(({ children,
    height,
    width,
    maxHeight,
    minHeight,
    maxWidth,
    minWidth,
    padding,
    margin,
    color,
    bgColor,
    border,
    radius,
    shadow,
    direction,
    justify,
    align,
    colGap,
    rowGap,
    basis,
    flex,
    shrink,
    grow,
    wrap,
    overflowX,
    overflowY,
    isBorderBox,
    pointer,
    ...RestProps
}: FlexProps, ref) => {

    return (
        <FlexContainer
            direction={direction}
            justify={justify}
            align={align}
            colGap={colGap}
            rowGap={rowGap}
            flex={flex}
            shrink={shrink}
            basis={basis}
            grow={grow}
            wrap={wrap}
            height={height}
            width={width}
            maxHeight={maxHeight}
            minHeight={minHeight}
            maxWidth={maxWidth}
            minWidth={minWidth}
            padding={padding}
            margin={margin}
            color={color}
            bgColor={bgColor}
            border={border}
            radius={radius}
            shadow={shadow}
            overflowX={overflowX}
            isBorderBox={isBorderBox}
            overflowY={overflowY}
            pointer={pointer}
            {...RestProps}
            ref={ref}
        >
            {children}
        </FlexContainer>
    )
})

const FlexContainer = styled.div<{
    height?: string,
    width?: string,
    maxHeight?: string,
    minHeight?: string,
    maxWidth?: string,
    minWidth?: string,
    padding?: string,
    margin?: string,
    color?: string,
    bgColor?: string,
    border?: string,
    radius?: string,
    shadow?: string,
    children: React.ReactNode;
    direction?: 'column' | 'row';
    justify?: string;
    align?: string;
    colGap?: string;
    rowGap?: string;
    flex?: number;
    shrink?: number;
    basis?: string;
    grow?: number;
    wrap?: 'nowrap' | 'wrap' | 'wrap-reverse';
    overflowX?: 'visible' | 'hidden' | 'clip' | 'scroll' | 'auto';
    overflowY?: 'visible' | 'hidden' | 'clip' | 'scroll' | 'auto';
    isBorderBox?: boolean
    pointer?: string;



}>(props => ({
    display: 'flex',
    height: props.height,
    width: props.width,
    maxHeight: props.maxHeight,
    minHeight: props.minHeight,
    maxWidth: props.maxWidth,
    minWidth: props.minWidth,
    padding: props.padding,
    margin: props.margin,
    color: props.color,
    backgroundColor: props.bgColor,
    border: props.border,
    borderRadius: props.radius,
    boxShadow: props.shadow,
    justifyContent: props.justify,
    alignItems: props.align,
    columnGap: props.colGap,
    gap: props.rowGap,
    flex: props.flex,
    flexShrink: props.shrink,
    flexGrow: props.grow,
    flexBasis: props.basis,
    flexWrap: props.wrap,
    overflowX: props.overflowX,
    overflowY: props.overflowY,
    boxSizing: 'border-box',
    pointer: props.pointer,
    flexDirection: props.direction,

}))

export default Flex 