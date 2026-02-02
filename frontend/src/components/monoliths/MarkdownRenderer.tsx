import React, { useEffect, useRef, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';
import 'katex/dist/katex.min.css';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism';

interface MarkdownRendererProps extends React.HTMLAttributes<HTMLDivElement> {
  content: string;
  highlightTerms?: string[];
}

const generateId = (text: string) => {
  return text
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-');
};



const MarkdownRenderer = React.memo(({
  content,
  highlightTerms = [],
  ...RestProps
}: MarkdownRendererProps) => {


  const containerRef = useRef<HTMLDivElement>(null);
  const itemRefs = useRef<HTMLElement[]>([]);
  const [activeIndex, setActiveIndex] = useState(-1)


  useEffect(() => {
    if (!containerRef.current || activeIndex === -1) return;

    const elements = containerRef.current.querySelectorAll('h1, h2, h3, h4, p, li, img, .table-container');
    const target = elements[activeIndex] as HTMLElement;

    if (target) {
      target.scrollIntoView({
        behavior: 'smooth',
        block: 'center'
      });
    }
  }, [activeIndex]);


  const highlightText = (text: string): React.ReactNode => {
    if (!highlightTerms.length) return text;

    const regex = new RegExp(
      `(${highlightTerms.map(term => term.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')).join('|')})`,
      'gi'
    );











    const parts = text.split(regex);

    return parts.map((part, index) => {
      const isHighlighted = highlightTerms.some(
        term => part.toLowerCase() === term.toLowerCase()
      );

      return isHighlighted ? (
        <mark
          key={index}
          style={{
            backgroundColor: '#fef08a',
            padding: '0em',
            borderRadius: '0.2em',
            fontWeight: 500
          }}
        >
          {part}
        </mark>
      ) : (
        <React.Fragment key={index}>{part}</React.Fragment>
      );
    });
  };

  const processChildren = (children: any): any => {
    if (typeof children === 'string') {
      return highlightText(children);
    }
    if (Array.isArray(children)) {
      return children.map((child, idx) =>
        React.isValidElement(child)
          ? React.cloneElement(child, { key: idx } as any)
          : typeof child === 'string'
            ? highlightText(child)
            : child
      );
    }
    return children;
  };


  const registerItem = (el: HTMLElement | null, index: number) => {
    if (el) {
      itemRefs.current[index] = el;
    }
  };


  const handleKeyDown = (e: React.KeyboardEvent<HTMLDivElement>) => {
    if (!containerRef.current) return;

    const elements = containerRef.current.querySelectorAll('h1, h2, h3, h4, p, li, img, .table-container');
    const max = elements.length - 1;

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setActiveIndex(prev => Math.min(prev + 1, max));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setActiveIndex(prev => Math.max(prev - 1, 0));
    }
  };


  const activeStyle = `
    .navigation-container > * { transition: all 0.2s; }
    .nav-active-element {
        background: #dbeafe !important;
        outline: 2px solid #3b82f6 !important;
        border-radius: 4px;
    }
    .table-container th, .table-container td {
        border-bottom: 1px solid #e5e7eb;
        border-right: 1px solid #e5e7eb;
    }
    .table-container th:last-child, .table-container td:last-child {
        border-right: none;
    }
    .table-container tr:last-child th, .table-container tr:last-child td {
        border-bottom: none;
    }
  `;

  useEffect(() => {
    if (!containerRef.current) return;
    const elements = containerRef.current.querySelectorAll('h1, h2, h3, h4, p, li, img, .table-container');

    elements.forEach((el, idx) => {
      if (idx === activeIndex) {
        el.classList.add('nav-active-element');
      } else {
        el.classList.remove('nav-active-element');
      }
    });
  }, [activeIndex]);



  const components = {

    h1: ({ children, ...props }: any) => {
      const text = React.Children.toArray(children).join('');
      const id = generateId(text);


      return (
        <h1
          id={id}

          style={{
            scrollMarginTop: '80px', // Pour éviter que l'ancre soit cachée sous un header fixe
            fontSize: '2em',
            fontWeight: 'bold',
            marginTop: '1.5em',
            marginBottom: '0.5em',
            lineHeight: '1.2'
          }}
          {...props}
        >
          {processChildren(children)}
        </h1>
      );
    },

    h2: ({ children, ...props }: any) => {
      const text = React.Children.toArray(children).join('');
      const id = generateId(text);

      return (
        <h2

          id={id}
          style={{
            scrollMarginTop: '80px',
            fontSize: '1.5em',
            fontWeight: 'bold',
            marginTop: '1.3em',
            marginBottom: '0.5em',
            lineHeight: '1.3'
          }}
          {...props}
        >
          {processChildren(children)}
        </h2>
      );
    },

    h3: ({ children, ...props }: any) => {
      const text = React.Children.toArray(children).join('');
      const id = generateId(text);

      return (
        <h3

          id={id}
          style={{
            scrollMarginTop: '80px',
            fontSize: '1.25em',
            fontWeight: 'bold',
            marginTop: '1.2em',
            marginBottom: '0.5em',
            lineHeight: '1.4'
          }}
          {...props}
        >
          {processChildren(children)}
        </h3>
      );
    },

    h4: ({ children, ...props }: any) => {
      const text = React.Children.toArray(children).join('');
      const id = generateId(text);

      return (
        <h4

          id={id}
          style={{
            scrollMarginTop: '80px',
            fontSize: '1.1em',
            fontWeight: 'bold',
            marginTop: '1em',
            marginBottom: '0.5em',
            lineHeight: '1.4'
          }}
          {...props}
        >
          {processChildren(children)}
        </h4>
      );
    },

    p: ({ children, ...props }: any) => {

      return (
        <p

          style={{
            marginTop: '1em',
            marginBottom: '1em',
            lineHeight: '1.7',
            fontSize: '1em'
          }}
          {...props}
        >
          {processChildren(children)}
        </p>
      )
    },

    li: ({ children, ...props }: any) => {

      return (
        <li
          style={{
            margin: "0.5em 0",
            lineHeight: "1.6",
            marginleft: "2.5rem",
          }}
          {...props}
        >
          <div
            style={{
              paddingLeft: "2.5rem",
            }}
          >
            {processChildren(children)}
          </div>
        </li>
      )


    },

    ul: ({ children, ...props }: any) => (
      <ul
        style={{
          marginTop: '1em',
          marginBottom: '1em',
          paddingLeft: '0',
          listStyleType: 'disc'
        }}
        {...props}
      >
        {children}
      </ul>
    ),

    ol: ({ children, ...props }: any) => (
      <ol
        style={{
          marginTop: '1em',
          marginLeft: '0',
          marginBottom: '1em',
          paddingLeft: '1.5em',
          listStyleType: 'decimal'
        }}
        {...props}
      >
        {children}
      </ol>
    ),

    code: ({ node, inline, className, children, ...props }: any) => {
      const match = /language-(\w+)/.exec(className || '');

      return !inline && match ? (
        <SyntaxHighlighter
          style={oneLight}
          language={match[1]}
          PreTag="div"
          customStyle={{
            borderRadius: '8px',
            fontSize: '0.9em',
            margin: '1em 0',
            padding: '1em'
          }}
          {...props}
        >
          {String(children).replace(/\n$/, '')}
        </SyntaxHighlighter>
      ) : (
        <code
          style={{
            backgroundColor: '#f3f4f6',
            padding: '0.2em 0.4em',
            borderRadius: '3px',
            fontSize: '0.9em',
            fontFamily: 'monospace'
          }}
          {...props}
        >
          {children}
        </code>
      );
    },

    pre: ({ children, ...props }: any) => (
      <pre
        style={{
          //backgroundColor: '#1e1e1e',
          borderRadius: '8px',
          padding: "1rem",
          overflowX: 'auto',
          margin: '1em 0'
        }}
        {...props}
      >
        {children}
      </pre>
    ),

    blockquote: ({ children, ...props }: any) => (
      <blockquote
        style={{
          borderLeft: '4px solid #3b82f6',
          paddingLeft: '1em',
          fontStyle: 'italic',
          margin: '1.5em 0',
          color: '#4b5563',
          backgroundColor: '#f9fafb',
          padding: '1em',
          borderRadius: '4px'
        }}
        {...props}
      >
        {processChildren(children)}
      </blockquote>
    ),

    strong: ({ children, ...props }: any) => (
      <strong style={{ fontWeight: 'bold' }} {...props}>
        {processChildren(children)}
      </strong>
    ),

    em: ({ children, ...props }: any) => (
      <em style={{ fontStyle: 'italic' }} {...props}>
        {processChildren(children)}
      </em>
    ),

    a: ({ href, children, ...props }: any) => {
      const isExternal = href?.startsWith('http');
      const isInternal = href?.startsWith('#');

      if (isInternal) {
        return (
          <a
            href={href}
            onClick={(e) => {
              e.preventDefault();
              const element = document.getElementById(href.slice(1));
              element?.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }}
            style={{
              color: '#2563eb',
              textDecoration: 'none',
              cursor: 'pointer'
            }}
            onMouseEnter={(e) => e.currentTarget.style.textDecoration = 'underline'}
            onMouseLeave={(e) => e.currentTarget.style.textDecoration = 'none'}
            {...props}
          >
            {children}
          </a>
        );
      }

      return (
        <a
          href={href}
          target={isExternal ? "_blank" : undefined}
          rel={isExternal ? "noopener noreferrer" : undefined}
          style={{
            color: '#2563eb',
            textDecoration: 'none'
          }}
          onMouseEnter={(e) => e.currentTarget.style.textDecoration = 'underline'}
          onMouseLeave={(e) => e.currentTarget.style.textDecoration = 'none'}
          {...props}
        >
          {children}
          {isExternal && <span style={{ marginLeft: '0.25em' }}>↗</span>}
        </a>
      );
    },

    img: ({ src, alt, ...props }: any) => {


      return (
        <img

          src={src}
          alt={alt}
          loading="lazy"
          style={{
            maxWidth: '100%',
            height: 'auto',
            borderRadius: '8px',
            boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
            margin: '1.5em 0',
            display: 'block'
          }}
          onError={(e) => {
            e.currentTarget.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="200" height="200"%3E%3Crect fill="%23ddd" width="200" height="200"/%3E%3Ctext fill="%23999" x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle"%3EImage unavailable%3C/text%3E%3C/svg%3E';
          }}
          {...props}
        />
      )

    },

    table: ({ children, ...props }: any) => {
      return (
        <div
          className="table-container"
          style={{
            width: "100%",
            maxWidth: "100%",
            overflowX: 'auto',
            boxSizing: "border-box",
            //border: ".5px solid lightgrey",
            borderRadius: "8px",
            border: '1px solid #e5e7eb',
            //borderTop: "none",
          }}

        >
          <table
            style={{
              borderSpacing: 0,
              borderCollapse: "collapse",
              display: "table",
              width: "100%",
              minWidth: "132px",
            }}
            {...props}
          >
            {children}
          </table>
        </div>
      )
    },

    thead: ({ children, ...props }: any) => (
      <thead
        style={{
          backgroundColor: '#f3f4f6',
          borderBottom: "1px solid #e5e7eb",
        }}
        {...props}
      >
        {children}
      </thead>
    ),

    tbody: ({ children, ...props }: any) => (
      <tbody {...props}>
        {children}
      </tbody>
    ),

    tr: ({ children, ...props }: any) => (
      <tr
        style={{
        }}
        {...props}
      >
        {children}
      </tr>
    ),

    th: ({ children, ...props }: any) => (
      <th
        style={{
          padding: '0.75em 1em',
          textAlign: 'left',
          fontWeight: 'bold',
          backgroundColor: '#f9fafb',
          whiteSpace: 'nowrap'
        }}
        {...props}
      >
        {processChildren(children)}
      </th>
    ),

    td: ({ children, ...props }: any) => (
      <td
        style={{
          padding: '0.75em 1em',
        }}
        {...props}
      >
        {processChildren(children)}
      </td>
    ),

    hr: (props: any) => (
      <hr
        style={{
          border: 'none',
          //borderTop: '1px solid #e5e7eb',
          margin: '2em 0'
        }}
        {...props}
      />
    ),

    sup: ({ children, ...props }: any) => (
      <sup
        style={{
          color: '#2563eb',
          cursor: 'pointer',
          fontSize: '0.8em'
        }}
        onClick={(e) => {
          const footnoteId = `fn-${children}`;
          document.getElementById(footnoteId)?.scrollIntoView({
            behavior: 'smooth',
            block: 'center'
          });
        }}
        onMouseEnter={(e) => e.currentTarget.style.textDecoration = 'underline'}
        onMouseLeave={(e) => e.currentTarget.style.textDecoration = 'none'}
        {...props}
      >
        {children}
      </sup>
    )
  };

  return (
    <>
      <style>{activeStyle}</style>
      <div
        ref={containerRef}
        {...RestProps}
        tabIndex={0}
        onKeyDown={handleKeyDown}
      //className="prose prose-slate max-w-none"
      >
        <ReactMarkdown
          remarkPlugins={[remarkGfm, remarkMath]}
          rehypePlugins={[rehypeKatex]}
          components={components}
        >
          {content}
        </ReactMarkdown>
      </div>


    </>

  );
});

export default MarkdownRenderer;
