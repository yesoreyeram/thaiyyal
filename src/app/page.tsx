"use client";
import React from "react";
import { useRouter } from "next/navigation";
import { AppNavBar } from "../components/AppNavBar";
import { HeroComponent } from "../components/HeroComponent";
import { OpenWorkflowModal } from "../components/OpenWorkflowModal";
import { Workflow } from "../types/workflow";

export default function HomePage() {
  const router = useRouter();
  const [showOpenModal, setShowOpenModal] = React.useState(false);

  const handleNewWorkflow = () => {
    // For static export, navigate to page-enhanced instead
    router.push('/page-enhanced');
  };

  const handleOpenWorkflow = (workflow: Workflow) => {
    // Save the workflow ID to localStorage and navigate
    localStorage.setItem('current_workflow_id', workflow.id);
    router.push('/page-enhanced');
  };

  return (
    <div className="h-screen flex flex-col bg-gray-950">
      {/* Application Nav Bar */}
      <AppNavBar
        onNewWorkflow={handleNewWorkflow}
        onOpenWorkflow={() => setShowOpenModal(true)}
      />

      {/* Hero Component */}
      <div className="flex-1 overflow-hidden">
        <HeroComponent onCreateWorkflow={handleNewWorkflow} />
      </div>

      {/* Open Workflow Modal */}
      <OpenWorkflowModal
        isOpen={showOpenModal}
        onClose={() => setShowOpenModal(false)}
        onSelect={handleOpenWorkflow}
      />
    </div>
  );
}
