.PHONY: run-k3d-with-keda
run-k3d-with-keda:
	@make -C ./keda run

.PHONY: busola-enable
busola-enable:
	@make -C ./busola busola-enable

.PHONY: run
run: run-k3d-with-keda busola-enable

.PHONY: stop
stop:
	@make -C ./keda stop
