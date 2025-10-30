"use client";
import React from "react";

interface HeroComponentProps {
  onCreateWorkflow: () => void;
}

export function HeroComponent({ onCreateWorkflow }: HeroComponentProps) {
  return (
    <div className="h-full w-full overflow-y-auto bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950">
      {/* Above the fold - Hero Section */}
      <div className="min-h-screen flex items-center justify-center px-8">
        <div className="text-center max-w-3xl">
          {/* Icon/Logo */}
          <div className="mb-8 flex justify-center">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full blur-2xl opacity-30 animate-pulse"></div>
              <div className="relative text-9xl">⚡</div>
            </div>
          </div>
          
          {/* Title */}
          <h1 className="text-6xl font-bold mb-6">
            <span className="bg-gradient-to-r from-blue-400 via-purple-500 to-pink-500 bg-clip-text text-transparent">
              Welcome to Thaiyyal
            </span>
          </h1>
          
          {/* Subtitle */}
          <p className="text-2xl text-gray-400 mb-12 leading-relaxed">
            Build powerful workflows visually with our enterprise-grade workflow builder.
          </p>
          
          {/* CTA Button */}
          <button
            onClick={onCreateWorkflow}
            className="group relative px-10 py-5 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white text-xl font-semibold rounded-xl transition-all shadow-2xl hover:shadow-blue-500/50 hover:scale-105 mb-6"
            aria-label="Create New Workflow"
          >
            <span className="relative z-10 flex items-center gap-3">
              <span className="text-3xl">✨</span>
              <span>Create New Workflow</span>
              <span className="text-2xl group-hover:translate-x-1 transition-transform">→</span>
            </span>
            <div className="absolute inset-0 bg-gradient-to-r from-blue-400 to-purple-400 rounded-xl blur opacity-0 group-hover:opacity-30 transition-opacity"></div>
          </button>
          
          {/* Keyboard Shortcut Hint */}
          <p className="text-sm text-gray-600">
            Or press <kbd className="px-2 py-1 bg-gray-800 border border-gray-700 rounded text-gray-400">Ctrl</kbd> + <kbd className="px-2 py-1 bg-gray-800 border border-gray-700 rounded text-gray-400">N</kbd> to create a new workflow
          </p>
        </div>
      </div>
      
      {/* Below the fold - Features Section */}
      <div className="min-h-screen flex items-center justify-center px-8 py-20">
        <div className="max-w-6xl w-full">
          {/* Section Title */}
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-white mb-4">
              Why Choose Thaiyyal?
            </h2>
            <p className="text-lg text-gray-500">
              Everything you need to build production-ready workflows
            </p>
          </div>
          
          {/* Features Grid */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
            <div className="p-8 bg-gray-800/50 rounded-xl border border-gray-700/50 backdrop-blur-sm hover:border-gray-600/50 transition-all">
              <div className="text-5xl mb-4">🎨</div>
              <h3 className="text-xl text-gray-200 font-semibold mb-3">Visual Editor</h3>
              <p className="text-gray-400">Intuitive drag & drop interface for building complex workflows without writing code.</p>
            </div>
            
            <div className="p-8 bg-gray-800/50 rounded-xl border border-gray-700/50 backdrop-blur-sm hover:border-gray-600/50 transition-all">
              <div className="text-5xl mb-4">⚙️</div>
              <h3 className="text-xl text-gray-200 font-semibold mb-3">Powerful Nodes</h3>
              <p className="text-gray-400">40+ specialized node types for operations, control flow, error handling, and more.</p>
            </div>
            
            <div className="p-8 bg-gray-800/50 rounded-xl border border-gray-700/50 backdrop-blur-sm hover:border-gray-600/50 transition-all">
              <div className="text-5xl mb-4">🚀</div>
              <h3 className="text-xl text-gray-200 font-semibold mb-3">Production Ready</h3>
              <p className="text-gray-400">Enterprise-grade quality with modern dark theme and professional UX.</p>
            </div>
          </div>
          
          {/* Additional Features */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="p-6 bg-gray-800/30 rounded-lg border border-gray-700/30">
              <div className="flex items-start gap-4">
                <div className="text-3xl">💾</div>
                <div>
                  <h4 className="text-lg text-gray-200 font-medium mb-2">Workflow Management</h4>
                  <p className="text-sm text-gray-400">Save, load, and manage multiple workflows with in-browser persistence.</p>
                </div>
              </div>
            </div>
            
            <div className="p-6 bg-gray-800/30 rounded-lg border border-gray-700/30">
              <div className="flex items-start gap-4">
                <div className="text-3xl">🔗</div>
                <div>
                  <h4 className="text-lg text-gray-200 font-medium mb-2">Node Connections</h4>
                  <p className="text-sm text-gray-400">Seamlessly connect nodes to define data flow and execution paths.</p>
                </div>
              </div>
            </div>
            
            <div className="p-6 bg-gray-800/30 rounded-lg border border-gray-700/30">
              <div className="flex items-start gap-4">
                <div className="text-3xl">📋</div>
                <div>
                  <h4 className="text-lg text-gray-200 font-medium mb-2">JSON Export</h4>
                  <p className="text-sm text-gray-400">Export workflows as JSON for version control and sharing.</p>
                </div>
              </div>
            </div>
            
            <div className="p-6 bg-gray-800/30 rounded-lg border border-gray-700/30">
              <div className="flex items-start gap-4">
                <div className="text-3xl">⌨️</div>
                <div>
                  <h4 className="text-lg text-gray-200 font-medium mb-2">Keyboard Shortcuts</h4>
                  <p className="text-sm text-gray-400">Boost productivity with keyboard shortcuts for common actions.</p>
                </div>
              </div>
            </div>
          </div>
          
          {/* Bottom CTA */}
          <div className="text-center mt-16">
            <button
              onClick={onCreateWorkflow}
              className="px-8 py-3 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white text-lg font-semibold rounded-lg transition-all shadow-lg hover:shadow-blue-500/30"
              aria-label="Get Started"
            >
              Get Started Now
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
