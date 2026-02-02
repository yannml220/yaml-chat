import React, { useState, useRef } from 'react';
import { Sparkles, ChevronLeft, ChevronRight } from 'lucide-react';

interface SuggestionCarouselProps {
  onSelect?: (suggestion: string) => void;
}

const SUGGESTIONS = [
  "Summarize this chapter",
  "Explain the historical context",
  "Find key character interactions",
  "Extract technical terminology",
  "Check for logical fallacies",
  "Compare with previous section",
  "Highlight core arguments",
  "Translate to simplified prose",
  "Generate a quiz from this content"
];

const SuggestionCarousel: React.FC<SuggestionCarouselProps> = ({ onSelect }) => {
  const [hoveredIdx, setHoveredIdx] = useState<number | null>(null);
  const [hoveredBtn, setHoveredBtn] = useState<'left' | 'right' | null>(null);
  const scrollRef = useRef<HTMLDivElement>(null);

  const handleScroll = (direction: 'left' | 'right') => {
    if (scrollRef.current) {
      const scrollAmount = 300; // Pixels to scroll
      scrollRef.current.scrollBy({
        left: direction === 'left' ? -scrollAmount : scrollAmount,
        behavior: 'smooth'
      });
    }
  };

  const wrapperStyle: React.CSSProperties = {
    width: '100%',
    position: 'relative',
    overflow: 'hidden',
    display: 'flex',
    alignItems: 'center',
  };

  const containerStyle: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    gap: '12px',
    padding: '16px 40px', // Horizontal padding to clear buttons
    overflowX: 'auto',
    scrollbarWidth: 'none', // Firefox
    msOverflowStyle: 'none', // IE/Edge
    scrollSnapType: 'x mandatory',
    WebkitOverflowScrolling: 'touch',
    width: '100%',
  };

  const navButtonStyle = (side: 'left' | 'right'): React.CSSProperties => ({
    position: 'absolute',
    [side]: side === 'left' ? '4px' : '4px',
    zIndex: 20,
    width: '32px',
    height: '32px',
    borderRadius: '50%',
    backgroundColor: 'rgba(255, 255, 255, 0.95)',
    border: '1px solid #e2e8f0',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
    color: '#64748b',
    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
    transition: 'all 0.2s ease',
    opacity: hoveredBtn === side ? 1 : 0.8,
    transform: hoveredBtn === side ? 'scale(1.1)' : 'scale(1)',
  });

  const itemBaseStyle: React.CSSProperties = {
    flex: '0 0 auto',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    padding: '8px 16px',
    borderRadius: '9999px',
    backgroundColor: '#ffffff',
    border: '1px solid #e2e8f0',
    fontSize: '13px',
    fontWeight: '600',
    color: '#475569',
    cursor: 'pointer',
    transition: 'all 0.2s ease',
    userSelect: 'none',
    boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
    scrollSnapAlign: 'start',
  };

  const getHoverStyle = (idx: number): React.CSSProperties => {
    if (hoveredIdx === idx) {
      return {
        backgroundColor: '#f5f3ff',
        borderColor: '#c4b5fd',
        color: '#4f46e5',
        transform: 'translateY(-1px)',
        boxShadow: '0 4px 6px -1px rgba(79, 70, 229, 0.1)',
      };
    }
    return {};
  };

  const iconStyle: React.CSSProperties = {
    color: '#6366f1',
    flexShrink: 0,
  };

  return (
    <div style={wrapperStyle}>
      {/* Global CSS for scrollbar hiding - as style tag because it's required for webkit browsers */}
      <style>
        {`
          .suggestion-container::-webkit-scrollbar {
            display: none;
          }
        `}
      </style>

      {/* Navigation Buttons */}
      <button
        style={navButtonStyle('left')}
        onMouseEnter={() => setHoveredBtn('left')}
        onMouseLeave={() => setHoveredBtn(null)}
        onClick={() => handleScroll('left')}
        aria-label="Scroll Left"
      >
        <ChevronLeft size={18} />
      </button>

      <div 
        ref={scrollRef}
        className="suggestion-container" 
        style={containerStyle}
      >
        {SUGGESTIONS.map((text, idx) => (
          <div
            key={idx}
            style={{ ...itemBaseStyle, ...getHoverStyle(idx) }}
            onMouseEnter={() => setHoveredIdx(idx)}
            onMouseLeave={() => setHoveredIdx(null)}
            onClick={() => onSelect?.(text)}
          >
            <Sparkles size={14} style={iconStyle} />
            <span>{text}</span>
          </div>
        ))}
      </div>

      <button
        style={navButtonStyle('right')}
        onMouseEnter={() => setHoveredBtn('right')}
        onMouseLeave={() => setHoveredBtn(null)}
        onClick={() => handleScroll('right')}
        aria-label="Scroll Right"
      >
        <ChevronRight size={18} />
      </button>
    </div>
  );
};

export default SuggestionCarousel;

