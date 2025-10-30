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
    // Generate a new workflow ID
    const workflowId = `workflow-${Date.now()}`;
    router.push(`/workflow/${workflowId}`);
  };

  const handleOpenWorkflow = (workflow: Workflow) => {
    router.push(`/workflow/${workflow.id}`);
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
