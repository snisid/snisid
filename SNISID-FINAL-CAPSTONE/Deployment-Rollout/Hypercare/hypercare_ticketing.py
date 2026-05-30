#!/usr/bin/env python3
"""
SNISID Hypercare Ticketing and Issue Resolution Simulator
Simulates automatic ticket routing, SLA calculations, and incident escalation
to stabilize regional operations during the critical 30-day Post-GoLive window.
"""

import json
from datetime import datetime, timedelta

class HypercareIncident:
    def __init__(self, ticket_id, department, category, severity, description):
        self.ticket_id = ticket_id
        self.department = department
        self.category = category # 'Hardware', 'Network', 'Identity_Conflict', 'Software'
        self.severity = severity # 'S1', 'S2', 'S3', 'S4'
        self.description = description
        self.created_at = datetime.now()
        self.status = "Open"
        self.assigned_tier = None
        self.sla_deadline = None
        self.resolution_notes = ""

    def route_and_assign_sla(self):
        """
        Routes the ticket to the correct Support Tier and calculates resolution SLA.
        """
        # Assign Tier based on Severity & Category
        if self.severity == "S1" or self.category == "Software":
            self.assigned_tier = "Tier 3 - Core Engineering"
        elif self.severity == "S2" or self.category == "Identity_Conflict":
            self.assigned_tier = "Tier 2 - Regional Experts"
        else:
            self.assigned_tier = "Tier 1 - Helpdesk"

        # Calculate SLA Deadline based on severity
        sla_hours_map = {
            "S1": 0.5, # 30 minutes
            "S2": 2.0, # 2 hours
            "S3": 6.0, # 6 hours
            "S4": 24.0 # 24 hours
        }
        hours_to_add = sla_hours_map.get(self.severity, 24.0)
        self.sla_deadline = self.created_at + timedelta(hours=hours_to_add)

    def resolve(self, notes):
        self.status = "Resolved"
        self.resolution_notes = notes

    def to_dict(self):
        return {
            "ticket_id": self.ticket_id,
            "department": self.department,
            "category": self.category,
            "severity": self.severity,
            "description": self.description,
            "created_at": self.created_at.isoformat(),
            "status": self.status,
            "assigned_tier": self.assigned_tier,
            "sla_deadline": self.sla_deadline.isoformat() if self.sla_deadline else None,
            "resolution_notes": self.resolution_notes
        }

class HypercareTicketingSystem:
    def __init__(self):
        self.tickets = {}
        self.ticket_counter = 0

    def create_ticket(self, department, category, severity, description):
        self.ticket_counter += 1
        ticket_id = f"HT-TICKET-{1000 + self.ticket_counter}"
        
        incident = HypercareIncident(ticket_id, department, category, severity, description)
        incident.route_and_assign_sla()
        
        self.tickets[ticket_id] = incident
        print(f"[Hypercare Support] Ticket {ticket_id} created for {department}! Status: Open -> Routed to {incident.assigned_tier} (SLA: {severity})")
        return incident

    def get_summary(self):
        open_count = sum(1 for t in self.tickets.values() if t.status == "Open")
        resolved_count = sum(1 for t in self.tickets.values() if t.status == "Resolved")
        
        tier_distribution = {}
        for t in self.tickets.values():
            tier_distribution[t.assigned_tier] = tier_distribution.get(t.assigned_tier, 0) + 1
            
        return {
            "total_tickets": len(self.tickets),
            "open_tickets": open_count,
            "resolved_tickets": resolved_count,
            "tier_distribution": tier_distribution
        }

if __name__ == "__main__":
    system = HypercareTicketingSystem()
    
    print("="*70)
    print("             SNISID HYPERCARE POST-GOLIVE TICKETING SIMULATION")
    print("="*70)
    
    # Simulate a variety of incidents during week 1 of deployment
    t1 = system.create_ticket(
        department="Ouest",
        category="Network",
        severity="S1",
        description="Fibre link cut between main DC and Delmas BLC site. Solaired-powered Starlink fails to initialize."
    )
    
    t2 = system.create_ticket(
        department="Nord",
        category="Identity_Conflict",
        severity="S2",
        description="Dual-entry validation fails on 4 files. System blocks enrollment for matching CIN."
    )
    
    t3 = system.create_ticket(
        department="Sud",
        category="Hardware",
        severity="S3",
        description="Biometric USB camera lens got scratched during transport. Image quality score below OACI threshold."
    )
    
    t4 = system.create_ticket(
        department="Grand'Anse",
        category="Software",
        severity="S2",
        description="Local Edge Node database integrity error during delayed sync. Corrupted WAL file."
    )
    
    print("\n--- SIMULATING RESOLUTION PROCESS ---")
    # Resolve the critical S1 fibre link cut ticket (Tier 3 resolved it via Starlink backup)
    system.tickets[t1.ticket_id].resolve(
        "Fibre link remained down, but Starlink terminal was manually reset and secondary satellite antenna was aligned. Connection stable."
    )
    print(f"[Hypercare Support] Resolved ticket {t1.ticket_id}: {system.tickets[t1.ticket_id].resolution_notes}")
    
    # Resolve the hardware lens ticket
    system.tickets[t3.ticket_id].resolve(
        "Field Strike team replaced the physical camera with a backup unit from the department hub. Standard OACI validation tests passed."
    )
    print(f"[Hypercare Support] Resolved ticket {t3.ticket_id}: {system.tickets[t3.ticket_id].resolution_notes}")

    # Generate final summary of Hypercare
    summary = system.get_summary()
    print("\n" + "="*70)
    print("                         HYPERCARE STATUS REPORT")
    print("="*70)
    print(f"Total Tickets Registered : {summary['total_tickets']}")
    print(f"Open / Active Tickets    : {summary['open_tickets']}")
    print(f"Resolved Tickets         : {summary['resolved_tickets']}")
    print("\nSupport Routing Distribution:")
    for tier, count in summary['tier_distribution'].items():
        print(f"  * {tier}: {count} ticket(s)")
    print("="*70)
    
    # Save tickets to report file
    serialized_tickets = {tid: t.to_dict() for tid, t in system.tickets.items()}
    with open("Deployment-Rollout/Hypercare/hypercare_tickets_report.json", "w") as f:
        json.dump(serialized_tickets, f, indent=4)
