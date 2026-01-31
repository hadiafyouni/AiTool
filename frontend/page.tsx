"use client";
import { useState } from "react";
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';

export default function CodingAssistant() {
  const [prompt, setPrompt] = useState("");
  const [answer, setAnswer] = useState("");
  const [loading, setLoading] = useState(false);

  const askAI = async () => {
    if (!prompt.trim()) return;
    setLoading(true);
    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/ask`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ prompt: prompt }),
      });
      const data = await res.json();
      setAnswer(data.answer);
    } catch (err) {
      setAnswer("Error connecting to the Go backend.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-slate-900 text-white p-8">
      <div className="max-w-3xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold border-b border-slate-700 pb-4">AI Coding Assistant</h1>

        <div className="flex flex-col gap-2">
          <textarea
            className="w-full p-4 bg-slate-800 rounded-lg border border-slate-700 outline-none h-32"
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
          />
          <button onClick={askAI} disabled={loading} className="bg-blue-600 p-3 rounded-lg">
            {loading ? "Thinking..." : "Generate Code"}
          </button>
        </div>

        {answer && (
          <div className="mt-8 p-6 bg-slate-800 rounded-lg border border-slate-700">
            <h2 className="text-sm font-semibold text-blue-400 mb-2">Response:</h2>

            <div className="prose prose-invert max-w-none">
              <ReactMarkdown
                components={{
                  code({ node, inline, className, children, ...props }: any) {
                    const match = /language-(\w+)/.exec(className || '');
                    return !inline && match ? (
                      <SyntaxHighlighter
                        style={vscDarkPlus}
                        language={match[1]}
                        PreTag="div"
                        {...props}
                      >
                        {String(children).replace(/\n$/, '')}
                      </SyntaxHighlighter>
                    ) : (
                      <code className={className} {...props}>
                        {children}
                      </code>
                    );
                  },
                }}
              >
                {answer}
              </ReactMarkdown>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}