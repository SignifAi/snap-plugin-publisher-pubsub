#!/bin/sh
snaptel task stop df6610aa-5ea8-4aee-ae43-d9d477bbafbc

snaptel plugin unload publisher pubsub-publisher 1
snaptel plugin load snap-plugin-publisher-pubsub

snaptel task enable e7f9f2f0-1f8e-42ab-b9f8-b9d1e88fafce
snaptel task start e7f9f2f0-1f8e-42ab-b9f8-b9d1e88fafce
